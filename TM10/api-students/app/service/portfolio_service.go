package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

// PortfolioService memakai DUA database: data mahasiswa dibaca dari
// PostgreSQL, portofolio disimpan di MongoDB.
type PortfolioService struct {
	portfolios repository.PortfolioRepository
	students   repository.StudentRepository
}

func NewPortfolioService(portfolios repository.PortfolioRepository, students repository.StudentRepository) *PortfolioService {
	return &PortfolioService{portfolios: portfolios, students: students}
}

func (s *PortfolioService) ensureStudentExists(c *fiber.Ctx, studentID int) error {
	if _, err := s.students.FindByID(c.Context(), studentID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("mahasiswa tidak ditemukan")
		}
		return err
	}
	return nil
}

// ---------- GET /students/:id/portfolio ----------
func (s *PortfolioService) Get(c *fiber.Ctx) error {
	studentID, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}
	if err := s.ensureStudentExists(c, studentID); err != nil {
		return err
	}

	portfolio, err := s.portfolios.GetByStudentID(c.Context(), studentID)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.OK(c, "berhasil mengambil portofolio", portfolio)
}

// ---------- PUT /students/:id/portfolio ----------
func (s *PortfolioService) Put(c *fiber.Ctx) error {
	studentID, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	var req model.PutPortfolioRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}
	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}

	if err := s.ensureStudentExists(c, studentID); err != nil {
		return err
	}

	portfolio, err := s.portfolios.Upsert(c.Context(), studentID, req)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.OK(c, "portofolio berhasil disimpan", portfolio)
}

// ---------- POST /students/:id/portfolio/projects ----------
func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
	studentID, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	var req model.ProjectInput
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}
	req.Title = strings.TrimSpace(req.Title)

	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}

	if err := s.ensureStudentExists(c, studentID); err != nil {
		return err
	}

	portfolio, err := s.portfolios.AddProject(c.Context(), studentID, model.Project{
		Title: req.Title, Technologies: req.Technologies, URL: req.URL, Year: req.Year,
	})
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Created(c, "proyek berhasil ditambahkan", portfolio)
}

// ---------- DELETE /students/:id/portfolio/projects/:projectId ----------
func (s *PortfolioService) RemoveProject(c *fiber.Ctx) error {
	studentID, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID mahasiswa tidak valid")
	}

	if err := s.ensureStudentExists(c, studentID); err != nil {
		return err
	}

	portfolio, err := s.portfolios.RemoveProject(c.Context(), studentID, c.Params("projectId"))
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.OK(c, "proyek berhasil dihapus", portfolio)
}

// ---------- GET /portfolios/stats/technologies ----------
func (s *PortfolioService) TechnologyStats(c *fiber.Ctx) error {
	stats, err := s.portfolios.TechnologyStats(c.Context())
	if err != nil {
		return err
	}
	return helper.OK(c, "berhasil mengambil statistik teknologi", stats)
}

func translatePortfolioError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("portofolio tidak ditemukan")
	case errors.Is(err, repository.ErrInvalidProjectID):
		return helper.BadRequest("id proyek tidak valid")
	case errors.Is(err, repository.ErrProjectLimit):
		return helper.Conflict("jumlah proyek sudah mencapai batas")
	default:
		return err
	}
}
