package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"api-students/app/handler"
	"api-students/app/service"
	"api-students/helper"
	"api-students/route"
)

func NewApp(studentService *service.StudentService, supervisionService *service.SupervisionService, portfolioService *service.PortfolioService, authHandler *handler.AuthHandler, jwtManager *helper.JWTManager, allowedOrigins string) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: helper.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
	}))
	app.Use(func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = helper.GenerateRandomString(16)
		}
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	})
	app.Use(LoggerMiddleware())

	route.RegisterRoutes(app, studentService, supervisionService, portfolioService, authHandler, jwtManager)

	return app
}
