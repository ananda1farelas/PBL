package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihan-fiber/app/handler"
	"latihan-fiber/app/repository"
	"latihan-fiber/app/service"
	"latihan-fiber/helper"
	"latihan-fiber/route"
)

func NewApp(pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Configuration / Environment Variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-change-in-production"
	}

	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	// 1. Inisialisasi Helper & Repository
	jwtManager := helper.NewJWTManager(jwtSecret, "latihan-fiber", accessTTL)
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)

	// 2. Inisialisasi Service
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, refreshTTL)

	// 3. Inisialisasi Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Registrasi Route (oper dependency lengkap)
	route.Register(app, pool, userHandler, authHandler, jwtManager)

	return app
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	slog.Error("global error handler", "status", code, "error", err.Error())
	return helper.Fail(c, code, err.Error())
}
