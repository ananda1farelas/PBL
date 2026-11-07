package model

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"` // Dicegah dari client (selalu diset server jadi 'user')
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthUser struct {
	ID       int
	Username string
	Role     string
}
