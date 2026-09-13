package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/app/handler" // Pastikan mengimpor package handler
	"latihan-fiber/helper"
	"latihan-fiber/middleware"
)

func Register(app *fiber.App, pool *pgxpool.Pool, userHandler *handler.UserHandler) {
	api := app.Group("/api/v1")

	api.Get("/health", healthCheck(pool))

	users := api.Group("/users", middleware.RequireJSON)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
	users.Post("/", userHandler.Create)
	users.Put("/:id", userHandler.Replace)
	users.Patch("/:id", userHandler.Patch)
	users.Delete("/:id", userHandler.Delete)
}

func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}

		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	}
}
