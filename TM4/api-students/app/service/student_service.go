package service

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type StudentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) List(c *fiber.Ctx) error {
	q := helper.ParseListQuery(c)
	students, total, err := s.repo.FindAll(c.Context(), q)
	if err != nil {
		return err
	}

	totalPages := CountTotalPages(total, q.Limit)
	meta := &model.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return helper.OKWithMeta(c, "berhasil mengambil daftar mahasiswa", students, meta)
}

func (s *StudentService) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, err.Error())
		}
		return err
	}

	return helper.OK(c, "berhasil mengambil detail mahasiswa", student)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format payload JSON tidak valid")
	}

	if errs := ValidateCreateStudent(req); len(errs) > 0 {
		return helper.ValidationError(c, "validasi gagal", errs)
	}

	student := model.Student{
		NIM:   req.NIM,
		Name:  req.Name,
		Email: req.Email,
	}

	created, err := s.repo.Create(c.Context(), student)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, err.Error())
		}
		return err
	}

	return helper.Created(c, "mahasiswa berhasil ditambahkan", created)
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format payload JSON tidak valid")
	}

	if errs := ValidateReplaceStudent(req); len(errs) > 0 {
		return helper.ValidationError(c, "validasi gagal", errs)
	}

	student := model.Student{
		ID:    id,
		NIM:   req.NIM,
		Name:  req.Name,
		Email: req.Email,
	}

	updated, err := s.repo.Update(c.Context(), student)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, err.Error())
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, err.Error())
		}
		return err
	}

	return helper.OK(c, "data mahasiswa berhasil diperbarui seluruhnya", updated)
}

func (s *StudentService) UpdatePartial(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format payload JSON tidak valid")
	}

	current, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, err.Error())
		}
		return err
	}

	updatedStudent, errs := ApplyPatchStudent(current, req)
	if len(errs) > 0 {
		return helper.ValidationError(c, "validasi gagal", errs)
	}

	result, err := s.repo.Update(c.Context(), updatedStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(c, fiber.StatusConflict, err.Error())
		}
		return err
	}

	return helper.OK(c, "data mahasiswa berhasil diperbarui sebagian", result)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	err = s.repo.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, err.Error())
		}
		return err
	}

	return helper.OK(c, "mahasiswa berhasil dihapus", nil)
}
