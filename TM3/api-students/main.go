package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"api-students/app/repository"
	"api-students/config"
	"api-students/database"
)

var metodeBerbody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// Middleware penolak request jika Content-Type bukan application/json (Status 415)
func requireJSON(c *fiber.Ctx) error {
	if metodeBerbody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	// 1. Memuat berkas .env ke environment proses
	config.LoadEnv()

	// 2. Membuka Connection Pool ke PostgreSQL
	pool, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer pool.Close()

	// 3. Perakitan Dependencies (Injeksi: Pool -> Repository -> Handler)
	studentRepo := repository.NewStudentRepository(pool)
	studentHandler := NewStudentHandler(studentRepo)

	// 4. Inisialisasi Fiber App
	app := fiber.New(fiber.Config{
		AppName: "API Students - Praktikum Backend Lanjut",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			pesan := "terjadi kesalahan pada server"

			if e, ok := err.(*fiber.Error); ok {
				status = e.Code
				pesan = e.Message
			}
			return fail(c, status, pesan)
		},
	})

	// Middleware Global
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())

	// Endpoint Root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("API Students Ready!")
	})

	api := app.Group("/api/v1")

	// Endpoint Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}
		return ok(c, "server dan database berjalan", nil)
	})

	// Endpoint Group Students
	s := api.Group("/students", requireJSON)
	s.Get("/", studentHandler.List)
	s.Get("/:id", studentHandler.GetByID)
	s.Post("/", studentHandler.Create)
	s.Put("/:id", studentHandler.Update)
	s.Patch("/:id", studentHandler.Patch)
	s.Delete("/:id", studentHandler.Delete)

	// Fallback Route (Status 404)
	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	// Port diambil secara dinamis dari file .env
	port := config.GetEnv("APP_PORT", "3000")
	log.Printf("Server berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}
