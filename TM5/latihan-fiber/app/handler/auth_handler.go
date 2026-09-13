package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/app/model"
	"latihan-fiber/app/service"
	"latihan-fiber/helper"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	user, err := h.authService.Register(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateUsername) || errors.Is(err, service.ErrDuplicateEmail) {
			return helper.Fail(c, fiber.StatusConflict, err.Error())
		}
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}

	return helper.Success(c, fiber.StatusCreated, "registrasi berhasil", user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	tokenPair, err := h.authService.Login(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return helper.Fail(c, fiber.StatusUnauthorized, err.Error())
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "terjadi kesalahan pada server")
	}

	return helper.Success(c, fiber.StatusOK, "login berhasil", tokenPair)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	tokenPair, err := h.authService.Refresh(c.UserContext(), req)
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, err.Error())
	}

	return helper.Success(c, fiber.StatusOK, "token berhasil diperbarui", tokenPair)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.RefreshRequest
	_ = c.BodyParser(&req)

	if err := h.authService.Logout(c.UserContext(), req.RefreshToken); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal melakukan logout")
	}

	return helper.Success(c, fiber.StatusOK, "logout berhasil", nil)
}
