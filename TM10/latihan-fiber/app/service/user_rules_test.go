package service

import (
	"errors"
	"strings"

	"latihan-fiber/app/model"
)

var (
	ErrInvalidInput = errors.New("input tidak valid")
)

func ValidateCreateUser(req model.CreateUserRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username wajib diisi")
	}
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email wajib diisi")
	}
	if len(req.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	return nil
}
