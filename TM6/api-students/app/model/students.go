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

type CreateStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive bool   `json:"isActive"`
}

type ReplaceStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive bool   `json:"isActive"`
}

type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty"`
	Name     *string `json:"name,omitempty"`
	Grade    *string `json:"grade,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}
