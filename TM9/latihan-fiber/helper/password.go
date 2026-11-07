package helper

import (
	"strings"
	"unicode"
)

var commonPasswords = map[string]bool{
	"password123": true,
	"12345678":    true,
	"qwerty123":   true,
}

// passwordStrength mengembalikan pesan kesalahan jika password lemah,
// atau string kosong jika password sudah cukup kuat.
func passwordStrength(password string) string {
	if len(password) < 8 {
		return "minimal 8 karakter"
	}
	if commonPasswords[strings.ToLower(password)] {
		return "password terlalu umum"
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "harus memuat huruf dan angka"
	}

	return ""
}
