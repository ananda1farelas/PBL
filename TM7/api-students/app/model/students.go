package model

import "time"

type Student struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Request untuk POST /students.
type CreateStudentRequest struct {
	NIM      string `json:"nim" validate:"required,nim"`
	Name     string `json:"name" validate:"required,notblank,min=2,max=50"`
	Grade    string `json:"grade" validate:"required,notblank"`
	IsActive bool   `json:"isActive"`
}

// Request untuk PUT /students/:id.
type ReplaceStudentRequest struct {
	NIM      string `json:"nim" validate:"required,nim"`
	Name     string `json:"name" validate:"required,notblank,min=2,max=50"`
	Grade    string `json:"grade" validate:"required,notblank"`
	IsActive bool   `json:"isActive"`
}

// Request untuk PATCH /students/:id.
// Pointer membedakan field yang tidak dikirim dengan field bernilai kosong.
type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty" validate:"omitnil,nim"`
	Name     *string `json:"name,omitempty" validate:"omitnil,notblank,min=2,max=50"`
	Grade    *string `json:"grade,omitempty" validate:"omitnil,notblank"`
	IsActive *bool   `json:"isActive,omitempty"`
}
