package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/helper"
)

func RequireAuth(jwtManager *helper.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return helper.Fail(c, fiber.StatusUnauthorized, "header Authorization tidak ditemukan")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.Fail(c, fiber.StatusUnauthorized, "format Authorization header harus 'Bearer <token>'")
		}

		tokenString := parts[1]
		authUser, err := jwtManager.Parse(tokenString)
		if err != nil {
			return helper.Fail(c, fiber.StatusUnauthorized, err.Error())
		}

		c.Locals(helper.LocalsAuthUser, authUser)
		return c.Next()
	}
}
