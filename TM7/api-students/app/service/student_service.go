package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strconv"
	"strings"

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
	format, err := helper.NegotiateStudentFormat(c.Get("Accept"))
	if err != nil {
		return err
	}

	q, err := helper.ParseCursorQuery(c)
	if err != nil {
		return helper.BadRequest("cursor tidak valid")
	}

	students, err := s.repo.FindAfterCursor(c.Context(), q)
	if err != nil {
		return err
	}

	// Ambil limit+1. Baris ekstra hanya menjadi penanda halaman berikutnya.
	hasMore := len(students) > q.Limit
	if hasMore {
		students = students[:q.Limit]
	}

	meta := &model.CursorMeta{Limit: q.Limit, HasMore: hasMore}
	if hasMore && len(students) > 0 {
		last := students[len(students)-1]
		meta.NextCursor = helper.EncodeCursor(last.CreatedAt, last.ID)
	}

	if format == helper.StudentFormatCSV {
		return sendStudentsCSV(c, students, meta)
	}

	return helper.OKWithCursorMeta(c, "berhasil mengambil daftar mahasiswa", students, meta)
}

func sendStudentsCSV(c *fiber.Ctx, students []model.Student, meta *model.CursorMeta) error {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	if err := writer.Write([]string{"id", "nim", "name", "grade", "is_active", "created_at"}); err != nil {
		return err
	}

	for _, student := range students {
		record := []string{
			strconv.Itoa(student.ID),
			student.NIM,
			student.Name,
			student.Grade,
			strconv.FormatBool(student.IsActive),
			student.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="students.csv"`)
	c.Set("X-Limit", strconv.Itoa(meta.Limit))
	c.Set("X-Has-More", strconv.FormatBool(meta.HasMore))
	if meta.NextCursor != "" {
		c.Set("X-Next-Cursor", meta.NextCursor)
	}
	return c.Status(fiber.StatusOK).Send(buffer.Bytes())
}

func (s *StudentService) GetByID(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID tidak valid")
	}

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("mahasiswa tidak ditemukan")
		}
		return err
	}

	return helper.OK(c, "berhasil mengambil detail mahasiswa", student)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Grade = strings.TrimSpace(req.Grade)

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}

	student := model.Student{
		NIM: req.NIM, Name: req.Name, Grade: req.Grade, IsActive: req.IsActive,
	}

	created, err := s.repo.Create(c.Context(), student)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM mahasiswa sudah terdaftar")
		}
		return err
	}

	return helper.Created(c, "mahasiswa berhasil ditambahkan", created)
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID tidak valid")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Grade = strings.TrimSpace(req.Grade)

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}

	student := model.Student{
		ID: id, NIM: req.NIM, Name: req.Name, Grade: req.Grade, IsActive: req.IsActive,
	}

	updated, err := s.repo.Update(c.Context(), student)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("mahasiswa tidak ditemukan")
		}
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM mahasiswa sudah digunakan")
		}
		return err
	}

	return helper.OK(c, "data mahasiswa berhasil diperbarui seluruhnya", updated)
}

func (s *StudentService) UpdatePartial(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID tidak valid")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}

	current, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("mahasiswa tidak ditemukan")
		}
		return err
	}

	updatedStudent := ApplyPatchStudent(current, req)
	result, err := s.repo.Update(c.Context(), updatedStudent)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM mahasiswa sudah digunakan")
		}
		return err
	}

	return helper.OK(c, "data mahasiswa berhasil diperbarui sebagian", result)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID tidak valid")
	}

	if err := s.repo.Delete(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("mahasiswa tidak ditemukan")
		}
		return err
	}

	return helper.OK(c, "mahasiswa berhasil dihapus", nil)
}
