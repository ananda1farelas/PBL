package helper

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort = errors.New("password minimal 8 karakter")
	ErrPasswordNoUpper  = errors.New("password harus mengandung minimal 1 huruf besar")
	ErrPasswordNoLower  = errors.New("password harus mengandung minimal 1 huruf kecil")
	ErrPasswordNoNumber = errors.New("password harus mengandung minimal 1 angka")
)

// ValidatePassword adalah function murni tanpa ketergantungan pada Fiber Ctx
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}

	return nil
}
