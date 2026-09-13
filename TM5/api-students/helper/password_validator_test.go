package helper_test

import (
	"testing"

	"api-students/helper" // Sesuaikan dengan nama module di go.mod kamu
)

func TestValidatePassword_TooShort(t *testing.T) {
	err := helper.ValidatePassword("Ab1")
	if err == nil || err != helper.ErrPasswordTooShort {
		t.Errorf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestValidatePassword_NoUpper(t *testing.T) {
	err := helper.ValidatePassword("password123")
	if err == nil || err != helper.ErrPasswordNoUpper {
		t.Errorf("expected ErrPasswordNoUpper, got %v", err)
	}
}

func TestValidatePassword_Valid(t *testing.T) {
	err := helper.ValidatePassword("Rahasia123")
	if err != nil {
		t.Errorf("expected no error for valid password, got %v", err)
	}
}
