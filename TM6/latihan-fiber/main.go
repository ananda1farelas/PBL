package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"latihan-fiber/app/repository"
	"latihan-fiber/app/service"
	"latihan-fiber/config"
	"latihan-fiber/database"
	"latihan-fiber/helper"
	"latihan-fiber/route"
)

const minSecretLength = 32

func main() {
	// 1. Konfigurasi dan logger
	config.LoadEnv()
	logger := config.NewLogger()

	// 2. Validasi JWT_SECRET
	// Rahasia diperiksa sebelum server menyala.
	jwtSecret := config.GetEnv("JWT_SECRET", "")

	if len(jwtSecret) < minSecretLength {
		logger.Error(
			"JWT_SECRET tidak diisi atau terlalu pendek",
			slog.Int("minimal_karakter", minSecretLength),
		)
		os.Exit(1)
	}

	// 3. Database
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error(
			"gagal terhubung ke database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer pool.Close()

	// 4. JWT Manager
	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "praktikum-backend"),
		time.Duration(
			config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15),
		)*time.Minute,
	)

	// 5. Repository
	userRepository := repository.NewUserRepository(pool)
	tokenRepository := repository.NewTokenRepository(pool)
	roleRepository := repository.NewRoleRepository(pool)

	// 6. Memuat permission dari database
	// Pemetaan role ke permission dibaca sekali
	// saat aplikasi pertama kali dijalankan.
	rawPermissions, err := roleRepository.LoadPermissions(
		context.Background(),
	)
	if err != nil {
		logger.Error(
			"gagal memuat permission",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	permissions := helper.NewPermissionSet(rawPermissions)

	logger.Info(
		"permission dimuat",
		slog.Any("roles", permissions.KnownRoles()),
	)

	// 7. Service
	userService := service.NewUserService(
		userRepository,
		permissions,
	)

	authService := service.NewAuthService(
		userRepository,
		tokenRepository,
		jwtManager,
		permissions,
		time.Duration(
			config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7),
		)*24*time.Hour,
	)

	// 8. Aplikasi
	app := config.NewApp(
		logger,
		route.Dependencies{
			Pool:        pool,
			JWT:         jwtManager,
			Permissions: permissions,
			UserService: userService,
			AuthService: authService,
		},
	)

	// 9. Jalankan server
	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error(
				"server berhenti",
				slog.String("error", err.Error()),
			)
			os.Exit(1)
		}
	}()

	logger.Info(
		"server berjalan",
		slog.String("port", port),
	)

	// 10. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logger.Info(
		"sinyal berhenti diterima, menutup server",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error(
			"gagal menutup server dengan rapi",
			slog.String("error", err.Error()),
		)
	}

	logger.Info("server berhenti dengan rapi")
}
