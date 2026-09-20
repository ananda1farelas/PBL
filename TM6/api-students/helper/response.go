package helper

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
)

func OK(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func OKWithMeta(c *fiber.Ctx, message string, data any, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Fail(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(model.WebResponse{
		Success: false,
		Message: message,
	})
}

func ValidationError(c *fiber.Ctx, message string, errs map[string]string) error {
	return c.Status(fiber.StatusBadRequest).JSON(model.WebResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	})
}
