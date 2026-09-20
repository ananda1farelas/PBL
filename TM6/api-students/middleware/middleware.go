package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
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
