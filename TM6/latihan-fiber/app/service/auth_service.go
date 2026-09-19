package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"latihan-fiber/app/model"
	"latihan-fiber/app/repository"
	"latihan-fiber/helper"
)

var (
	ErrInvalidCredentials  = errors.New("username atau password salah")
	ErrDuplicateUsername   = errors.New("username sudah digunakan")
	ErrDuplicateEmail      = errors.New("email sudah digunakan")
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid atau sudah kedaluwarsa")
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.User, error)
	Login(ctx context.Context, req model.LoginRequest) (model.TokenPair, error)
	Refresh(ctx context.Context, req model.RefreshRequest) (model.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *helper.JWTManager
	refreshTTL time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager *helper.JWTManager,
	refreshTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return model.User{}, errors.New("username, email, dan password wajib diisi")
	}

	// Cek apakah username sudah dipakai
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser.ID != 0 {
		return model.User{}, ErrDuplicateUsername
	}

	// Hash password sebelum disimpan
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("hashing password: %w", err)
	}

	u := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	createdUser, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return model.User{}, fmt.Errorf("menyimpan user: %w", err)
	}

	return createdUser, nil
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (model.TokenPair, error) {
	req.Username = strings.TrimSpace(req.Username)

	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// Jalankan perbandingan dummy password untuk mencegah timing attack
		helper.VerifyDummyPassword(req.Password)
		return model.TokenPair{}, ErrInvalidCredentials
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return model.TokenPair{}, ErrInvalidCredentials
	}

	return s.generateTokenPair(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, req model.RefreshRequest) (model.TokenPair, error) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return model.TokenPair{}, ErrInvalidRefreshToken
	}

	tokenHash := helper.SHA256Hex(req.RefreshToken)
	token, err := s.tokenRepo.FindActive(ctx, tokenHash)
	if err != nil {
		return model.TokenPair{}, ErrInvalidRefreshToken
	}

	// Token lama langsung dicabut (Revoke-on-Use pattern)
	if err := s.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return model.TokenPair{}, fmt.Errorf("mencabut token lama: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return model.TokenPair{}, ErrInvalidRefreshToken
	}

	return s.generateTokenPair(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}

	tokenHash := helper.SHA256Hex(refreshToken)
	return s.tokenRepo.Revoke(ctx, tokenHash)
}

// Helper internal untuk membuat access token + refresh token pair
func (s *authService) generateTokenPair(ctx context.Context, u model.User) (model.TokenPair, error) {
	accessToken, err := s.jwtManager.GenerateAccess(u)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("membuat access token: %w", err)
	}

	rawRefreshToken, err := helper.RandomToken(32)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("membuat raw refresh token: %w", err)
	}

	refreshTokenHash := helper.SHA256Hex(rawRefreshToken)
	expiresAt := time.Now().Add(s.refreshTTL)

	rf := model.RefreshToken{
		UserID:    u.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: expiresAt,
	}

	if err := s.tokenRepo.Save(ctx, rf); err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtManager.AccessTTL().Seconds()),
	}, nil
}
