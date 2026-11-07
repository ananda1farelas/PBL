package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,nospace"`
	Role     string `json:"role" validate:"required,oneof=admin staff user"`
}

type ReplaceUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,nospace"`
	Role     string `json:"role" validate:"required,oneof=admin staff user"`
}

type PatchUserRequest struct {
	Username *string `json:"username,omitempty" validate:"min=3,max=30,alphanum"`
	Email    *string `json:"email,omitempty" validate:"email,max=120"`
	Password *string `json:"password,omitempty" validate:"min=8,max=72,nospace"`
	Role     *string `json:"role,omitempty" validate:"oneof=admin staff user"`
}

type AssignRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin staff user"`
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}
