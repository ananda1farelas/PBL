package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

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
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
}

type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager *helper.JWTManager
	refreshTTL time.Duration
	perms      *helper.PermissionSet
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager *helper.JWTManager,
	perms *helper.PermissionSet,
	refreshTTL time.Duration,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		perms:      perms,
		refreshTTL: refreshTTL,
	}
}

func (s *authService) Register(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"username, email, dan password wajib diisi",
		)
	}

	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser.ID != 0 {
		return helper.Fail(
			c,
			fiber.StatusConflict,
			ErrDuplicateUsername.Error(),
		)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal memproses password",
		)
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
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat user",
		)
	}

	return helper.Success(
		c,
		fiber.StatusCreated,
		"registrasi berhasil",
		createdUser,
	)
}

func (s *authService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	req.Username = strings.TrimSpace(req.Username)

	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		helper.VerifyDummyPassword(req.Password)

		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			ErrInvalidCredentials.Error(),
		)
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			ErrInvalidCredentials.Error(),
		)
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat token",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"login berhasil",
		tokenPair,
	)
}

func (s *authService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			ErrInvalidRefreshToken.Error(),
		)
	}

	tokenHash := helper.SHA256Hex(req.RefreshToken)

	token, err := s.tokenRepo.FindActive(ctx, tokenHash)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			ErrInvalidRefreshToken.Error(),
		)
	}

	if err := s.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mencabut token lama",
		)
	}

	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			ErrInvalidRefreshToken.Error(),
		)
	}

	tokenPair, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat token baru",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"token berhasil diperbarui",
		tokenPair,
	)
}

func (s *authService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
		)
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.NoContent(c)
	}

	tokenHash := helper.SHA256Hex(req.RefreshToken)

	if err := s.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal logout",
		)
	}

	return helper.NoContent(c)
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

func (s *authService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	user, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
		)
	}

	_ = ctx

	return helper.Success(
		c,
		fiber.StatusOK,
		"profil berhasil diambil",
		fiber.Map{
			"user":        user,
			"permissions": s.perms.PermissionsOf(user.Role),
		},
	)
}
