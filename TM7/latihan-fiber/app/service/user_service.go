package service

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/helper"
)

type UserService struct {
	repo  repository.UserRepository
	perms *helper.PermissionSet
}

func NewUserService(
	repo repository.UserRepository,
	perms *helper.PermissionSet,
) *UserService {
	return &UserService{
		repo:  repo,
		perms: perms,
	}
}

// ---------- LIST USERS ----------

func (s *UserService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	users, err := s.repo.List(ctx)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil data user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"berhasil mengambil daftar user",
		users,
	)
}

// ---------- FIND USER BY ID ----------

func (s *UserService) FindByID(
	ctx context.Context,
	id int,
) (model.User, error) {
	return s.repo.FindByID(ctx, id)
}

// ---------- CREATE USER ----------

func (s *UserService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	if err := ValidateCreate(req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"input tidak valid",
		)
	}

	u := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		IsActive: true,
	}

	user, err := s.repo.Create(ctx, u)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal membuat user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusCreated,
		"user berhasil dibuat",
		user,
	)
}

// ---------- REPLACE USER ----------

func (s *UserService) Replace(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	var req model.ReplaceUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	u := model.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.repo.Update(ctx, id, u)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal memperbarui user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user berhasil diperbarui",
		user,
	)
}

// ---------- PATCH USER ----------

func (s *UserService) Patch(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	var req model.PatchUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	current, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal mengambil data user",
		)
	}

	updated, errs := ApplyPatch(current, req)

	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	user, err := s.repo.Update(ctx, id, updated)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal memperbarui user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user berhasil diperbarui secara parsial",
		user,
	)
}

// ---------- GET /users/:id ----------

func (s *UserService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
		)
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	// Pemeriksaan hak akses dilakukan sebelum data diambil.
	if !CanAccessUser(
		current,
		id,
		s.perms,
		"user:read:any",
	) {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak berhak mengakses data user lain",
		)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal mengambil data user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user ditemukan",
		user,
	)
}

// ---------- PATCH /users/:id/role ----------

// Dijaga middleware dengan permission role:assign.
func (s *UserService) AssignRole(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
		)
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	var req model.AssignRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	if errs := ValidateAssignRole(
		current,
		id,
		req,
		s.perms,
	); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	result, err := s.repo.UpdateRole(
		ctx,
		id,
		strings.TrimSpace(req.Role),
	)
	if err != nil {
		return translateError(
			c,
			err,
			"gagal mengubah role user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"role user berhasil diubah",
		result,
	)
}

// ---------- DELETE /users/:id ----------

func (s *UserService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
		)
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	// Punya permission menghapus tidak berarti
	// boleh menghapus dirinya sendiri.
	if current.UserID == id {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak boleh menghapus akun sendiri",
		)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateError(
			c,
			err,
			"gagal menghapus user",
		)
	}

	return helper.NoContent(c)
}

func translateError(
	c *fiber.Ctx,
	err error,
	message string,
) error {
	if errors.Is(err, repository.ErrNotFound) {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"data tidak ditemukan",
		)
	}

	if errors.Is(err, repository.ErrDuplicate) {
		return helper.Fail(
			c,
			fiber.StatusConflict,
			"data sudah ada",
		)
	}

	return helper.Fail(
		c,
		fiber.StatusInternalServerError,
		message,
	)
}
