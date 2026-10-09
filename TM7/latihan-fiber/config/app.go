package config

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"latihan-fiber/helper"
	"latihan-fiber/middleware"
	"latihan-fiber/route"
)

// NewApp merakit aplikasi: membuat instance Fiber, memasang middleware,
// lalu mendaftarkan route. File ini adalah tempat seluruh bagian bertemu.
func NewApp(
	logger *slog.Logger,
	deps route.Dependencies,
) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: GetEnv(
			"APP_NAME",
			"Praktikum Backend Lanjut",
		),
		ErrorHandler: newErrorHandler(logger),
		BodyLimit:    1 * 1024 * 1024,
	})

	middleware.Register(app, logger)
	app.Use(middleware.RequireJSONAccept)

	route.Register(app, deps)

	app.Use(func(c *fiber.Ctx) error {
		return helper.NotFound("endpoint tidak ditemukan")
	})

	return app
}
