package handler

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/service"
	"api-students/helper"
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

	if err := h.authService.Register(c.UserContext(), req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}

	return helper.Created(c, "registrasi berhasil", nil)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	tokens, err := h.authService.Login(c.UserContext(), req)
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, err.Error())
	}

	return helper.OK(c, "login berhasil", tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req model.TokenRefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	tokens, err := h.authService.Refresh(c.UserContext(), req.RefreshToken)
	if err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, err.Error())
	}

	return helper.OK(c, "token berhasil di-refresh", tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.TokenRefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "format body request tidak valid")
	}

	_ = h.authService.Logout(c.UserContext(), req.RefreshToken)
	return helper.OK(c, "logout berhasil", nil)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	authUser, ok := c.Locals(helper.LocalsAuthUser).(*helper.AuthClaims)
	if !ok || authUser == nil {
		return helper.Fail(c, fiber.StatusUnauthorized, "akses ditolak")
	}

	profile, err := h.authService.GetProfile(c.UserContext(), authUser.Username)
	if err != nil {
		return helper.Fail(c, fiber.StatusNotFound, err.Error())
	}

	return helper.OK(c, "profil berhasil diambil", profile)
}
