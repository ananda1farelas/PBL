package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/helper"
)

func RequireJSON(c *fiber.Ctx) error {
	method := c.Method()
	if method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodPatch {
		contentType := c.Get("Content-Type")
		if !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
			return helper.UnsupportedMediaType("Content-Type harus application/json")
		}
	}
	return c.Next()
}

func JWTAuth(manager *helper.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "token akses tidak ditemukan",
			})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := manager.Parse(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "token akses tidak valid atau kedaluwarsa",
			})
		}

		c.Locals(helper.LocalsAuthUser, claims)
		return c.Next()
	}
}

// RequirePermission memeriksa permission berdasarkan role pada access token.
func RequirePermission(permissions *helper.PermissionSet, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(helper.LocalsAuthUser).(*helper.AuthClaims)
		if !ok || claims == nil {
			return helper.Unauthorized("akses ditolak")
		}
		if !permissions.Has(claims.Role, permission) {
			return helper.Forbidden("tidak memiliki permission: " + permission)
		}
		return c.Next()
	}
}
