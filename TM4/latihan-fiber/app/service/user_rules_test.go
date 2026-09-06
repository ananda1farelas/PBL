package service

import (
	"testing"

	"api-students/app/model" // Pastikan 'api-students' sesuai dengan nama di go.mod
)

func TestValidateCreate(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		req := model.CreateUserRequest{
			Name:  "Budi Santoso",
			Email: "budi@example.com",
		}
		errs := ValidateCreate(req)
		if len(errs) != 0 {
			t.Errorf("diharapkan tidak ada error, dapat: %v", errs)
		}
	})

	t.Run("empty fields", func(t *testing.T) {
		req := model.CreateUserRequest{
			Name:  "",
			Email: "   ",
		}
		errs := ValidateCreate(req)
		if len(errs) != 2 {
			t.Errorf("diharapkan 2 error (name dan email), dapat: %d", len(errs))
		}
	})
}

func TestCountTotalPages(t *testing.T) {
	tests := []struct {
		total    int
		limit    int
		expected int
	}{
		{total: 0, limit: 10, expected: 0},
		{total: 5, limit: 10, expected: 1},
		{total: 10, limit: 10, expected: 1},
		{total: 11, limit: 10, expected: 2},
		{total: 25, limit: 10, expected: 3},
		{total: 10, limit: 0, expected: 0},
	}

	for _, tt := range tests {
		got := CountTotalPages(tt.total, tt.limit)
		if got != tt.expected {
			t.Errorf("CountTotalPages(%d, %d) = %d; ingin %d", tt.total, tt.limit, got, tt.expected)
		}
	}
}
