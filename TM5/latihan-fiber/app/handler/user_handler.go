package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
	"latihan-fiber/app/service"
	"latihan-fiber/helper"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.userService.List(c.UserContext())
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal mengambil data user")
	}
	return helper.Success(c, fiber.StatusOK, "berhasil mengambil daftar user", users)
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	user, err := h.userService.FindByID(c.UserContext(), id)
	if err != nil {
		return helper.Fail(c, fiber.StatusNotFound, "user tidak ditemukan")
	}
	return helper.Success(c, fiber.StatusOK, "berhasil mengambil data user", user)
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	user, err := h.userService.Create(c.UserContext(), req)
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return helper.Success(c, fiber.StatusCreated, "user berhasil dibuat", user)
}

func (h *UserHandler) Replace(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var req model.ReplaceUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	user, err := h.userService.Update(c.UserContext(), id, req)
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "user berhasil diperbarui", user)
}

func (h *UserHandler) Patch(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var req model.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	user, err := h.userService.Patch(c.UserContext(), id, req)
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}
	return helper.Success(c, fiber.StatusOK, "user berhasil diperbarui secara parsial", user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	if err := h.userService.Delete(c.UserContext(), id); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal menghapus user")
	}
	return helper.Success(c, fiber.StatusOK, "user berhasil dihapus", nil)
}
