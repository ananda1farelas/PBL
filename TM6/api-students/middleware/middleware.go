package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/helper"
)

func RequireJSON(c *fiber.Ctx) error {
	method := c.Method()
	if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch {
		contentType := c.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(model.WebResponse{
				Success: false,
				Message: "Content-Type harus application/json",
			})
		}
	}
	return c.Next()
}

func JWTAuth(manager *helper.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(model.WebResponse{
				Success: false,
				Message: "token akses tidak ditemukan",
			})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := manager.Parse(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.WebResponse{
				Success: false,
				Message: "token akses tidak valid atau kedaluwarsa",
			})
		}

		c.Locals(helper.LocalsAuthUser, claims)
		return c.Next()
	}
}