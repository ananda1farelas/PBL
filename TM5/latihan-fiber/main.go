package main

import (
	"context"
	"log/slog"
	"os"

	"latihan-fiber/config"
	"latihan-fiber/database"
)

func main() {
	ctx := context.Background()

	// 1. Inisialisasi Database
	pool, err := database.NewPostgresPool(ctx)
	if err != nil {
		slog.Error("gagal terhubung ke database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 2. Inisialisasi Fiber App
	app := config.NewApp(pool)

	// 3. Jalankan Server
	slog.Info("server berjalan pada port :3000")
	if err := app.Listen(":3000"); err != nil {
		slog.Error("gagal menjalankan server", "error", err)
	}
}
