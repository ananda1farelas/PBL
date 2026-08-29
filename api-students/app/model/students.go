package model

import "time"

// Student adalah struct utama tabel database
type Student struct {
	ID        int64     `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Request DTO untuk POST (Tambah Mahasiswa)
type CreateStudentRequest struct {
	NIM      string  `json:"nim"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive *bool   `json:"is_active"`
}

// Request DTO untuk PUT/PATCH (Update Mahasiswa)
type UpdateStudentRequest struct {
	NIM      string  `json:"nim"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// WebResponse adalah pembungkus standar response API JSON
type WebResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Meta    *Meta             `json:"meta,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Meta menyimpan informasi pagination
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListQuery menampung parameter pagination, search, sort, dan filter
type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	SortBy   string
	Order    string
	IsActive *bool
}

// Offset menghitung offset database untuk klausa LIMIT & OFFSET
func (q ListQuery) Offset() int {
	if q.Page <= 1 {
		return 0
	}
	return (q.Page - 1) * q.Limit
}
