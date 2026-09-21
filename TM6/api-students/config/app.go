package config

import (
	"github.com/gofiber/fiber/v2"

	"api-students/app/handler"
	"api-students/app/service"
	"api-students/helper"
	"api-students/route"
)

func NewApp(studentService *service.StudentService, authHandler *handler.AuthHandler, jwtManager *helper.JWTManager) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(LoggerMiddleware())

	route.RegisterRoutes(app, studentService, authHandler, jwtManager)

	return app
}