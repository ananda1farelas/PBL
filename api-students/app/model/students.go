package model

import "time"

// Entity utama Student sesuai skema tabel database
type Student struct {
	ID        int64     `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// DTO untuk request pembuatan mahasiswa baru (POST)
type CreateStudentRequest struct {
	NIM      string  `json:"nim"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive *bool   `json:"is_active"`
}

// DTO untuk request pembaruan mahasiswa (PUT)
type UpdateStudentRequest struct {
	NIM      string  `json:"nim"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// Parameter query string untuk pencarian, pengurutan, dan pagination
type ListQuery struct {
	Search string
	SortBy string
	Order  string
	Limit  int
	Page   int
}

// Helper untuk menghitung offset query database
func (q ListQuery) Offset() int {
	if q.Page <= 1 {
		return 0
	}
	return (q.Page - 1) * q.Limit
}