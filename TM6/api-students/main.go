package main

import (
	"context"
	"fmt"
	"log"

	"api-students/app/handler"
	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
)

func main() {
	cfg := config.LoadConfig()

	// Inisialisasi Koneksi Database PostgreSQL
	db, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtManager := helper.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTTL)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, cfg.RefreshTTL)
	authHandler := handler.NewAuthHandler(authService)

	app := config.NewApp(studentService, authHandler, jwtManager)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}