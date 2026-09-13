package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) error
	Login(ctx context.Context, req model.LoginRequest) (model.TokenPairResponse, error)
	Refresh(ctx context.Context, refreshTokenStr string) (model.TokenPairResponse, error)
	Logout(ctx context.Context, refreshTokenStr string) error
	GetProfile(ctx context.Context, username string) (model.UserResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *helper.JWTManager
	refreshTTL time.Duration
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtManager *helper.JWTManager, refreshTTL time.Duration) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(ctx context.Context, req model.RegisterRequest) error {
	// Validasi kekuatan password murni
	if err := helper.ValidatePassword(req.Password); err != nil {
		return err
	}

	// Hashing password dengan bcrypt cost 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	// Paksa role selalu 'user' (keamanan poin 4 soal: mengirim role admin tetap jadi user)
	role := "user"

	err = s.userRepo.Create(ctx, req.Username, string(hashedPassword), role)
	if err != nil {
		return errors.New("username sudah terdaftar atau gagal menyimpan user")
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req model.LoginRequest) (model.TokenPairResponse, error) {
	// Ambil data user dari repository
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		// Pesan error disamakan agar aman dari user enumeration (poin 4 soal)
		return model.TokenPairResponse{}, errors.New("username atau password salah")
	}

	// Verifikasi password bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return model.TokenPairResponse{}, errors.New("username atau password salah")
	}

	// Generate Access Token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.Username, user.Role)
	if err != nil {
		return model.TokenPairResponse{}, err
	}

	// Generate Refresh Token acak
	rawRefreshToken := helper.GenerateRandomString(32)
	tokenHash := s.hashToken(rawRefreshToken)
	expiresAt := time.Now().Add(s.refreshTTL)

	// Simpan hash refresh token ke database
	if err := s.tokenRepo.Save(ctx, user.Username, tokenHash, expiresAt); err != nil {
		return model.TokenPairResponse{}, err
	}

	return model.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshTokenStr string) (model.TokenPairResponse, error) {
	tokenHash := s.hashToken(refreshTokenStr)

	// Cari dan validasi refresh token di DB
	username, err := s.tokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return model.TokenPairResponse{}, errors.New("refresh token tidak valid atau sudah kedaluwarsa")
	}

	// Rotasi: Hapus refresh token lama agar tidak bisa dipakai ulang (poin 4 soal)
	_ = s.tokenRepo.DeleteByHash(ctx, tokenHash)

	// Ambil data user untuk role terbaru
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return model.TokenPairResponse{}, errors.New("user tidak ditemukan")
	}

	// Buat token pair baru
	accessToken, err := s.jwtManager.GenerateAccessToken(user.Username, user.Role)
	if err != nil {
		return model.TokenPairResponse{}, err
	}

	newRawRefresh := helper.GenerateRandomString(32)
	newHash := s.hashToken(newRawRefresh)
	expiresAt := time.Now().Add(s.refreshTTL)

	if err := s.tokenRepo.Save(ctx, user.Username, newHash, expiresAt); err != nil {
		return model.TokenPairResponse{}, err
	}

	return model.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: newRawRefresh,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshTokenStr string) error {
	tokenHash := s.hashToken(refreshTokenStr)
	return s.tokenRepo.DeleteByHash(ctx, tokenHash)
}

func (s *authService) GetProfile(ctx context.Context, username string) (model.UserResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return model.UserResponse{}, errors.New("user tidak ditemukan")
	}

	return model.UserResponse{
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *authService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
