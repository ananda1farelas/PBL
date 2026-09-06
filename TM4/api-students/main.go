package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
)

func main() {
	cfg := config.LoadConfig()

	// Inisialisasi Koneksi Database MySQL
	dsn := "postgres://postgres:postgres@127.0.0.1:5432/db_students"
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Gagal membuat koneksi database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal ping database: %v", err)
	}

	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)

	app := config.NewApp(studentService)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
