package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"api-students/app/handler"
	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	// MongoDB — portofolio mahasiswa
	mongoClient, err := database.NewMongoClient(context.Background(), cfg.MongoURI)
	if err != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoClient.Disconnect(ctx)
	}()
	mongoDB := mongoClient.Database(cfg.MongoDatabase)

	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)

	txManager := repository.NewTxManager(db)
	supervisorRepo := repository.NewSupervisorRepository(db)
	supervisionRepo := repository.NewSupervisionRepository(db)
	supervisionService := service.NewSupervisionService(txManager, supervisorRepo, supervisionRepo)

	portfolioRepo := repository.NewPortfolioRepository(mongoDB)
	if err := portfolioRepo.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("Gagal membuat index MongoDB: %v", err)
	}
	portfolioService := service.NewPortfolioService(portfolioRepo, studentRepo)

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtManager := helper.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.AccessTTL)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtManager, cfg.RefreshTTL)
	authHandler := handler.NewAuthHandler(authService)

	app := config.NewApp(studentService, supervisionService, portfolioService, authHandler, jwtManager, cfg.AllowedOrigins)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
