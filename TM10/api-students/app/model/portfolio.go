package model

import "time"

// Portfolio adalah bentuk portofolio yang dikenal seluruh aplikasi — tidak
// mengandung tipe apa pun dari driver MongoDB.
type Portfolio struct {
	ID           string        `json:"id"`
	StudentID    int           `json:"student_id"`
	Headline     string        `json:"headline"`
	Bio          string        `json:"bio"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Project struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Technologies []string `json:"technologies"`
	URL          string   `json:"url"`
	Year         int      `json:"year"`
}

type Certificate struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Date   string `json:"date"`
}

// PutPortfolioRequest dipakai untuk PUT /students/:id/portfolio — membuat
// atau mengganti seluruh portofolio sekaligus.
type PutPortfolioRequest struct {
	Headline     string             `json:"headline" validate:"max=150"`
	Bio          string             `json:"bio" validate:"max=1000"`
	Projects     []ProjectInput     `json:"projects" validate:"max=20,dive"`
	Certificates []CertificateInput `json:"certificates" validate:"max=50,dive"`
}

// ProjectInput dipakai baik di dalam PutPortfolioRequest maupun sebagai body
// POST /students/:id/portfolio/projects.
type ProjectInput struct {
	Title        string   `json:"title" validate:"required,notblank,max=150"`
	Technologies []string `json:"technologies" validate:"max=20,dive,notblank,max=50"`
	URL          string   `json:"url" validate:"omitempty,url"`
	Year         int      `json:"year" validate:"omitempty,min=2000,max=2100"`
}

type CertificateInput struct {
	Name   string `json:"name" validate:"required,notblank,max=150"`
	Issuer string `json:"issuer" validate:"required,notblank,max=150"`
	Date   string `json:"date" validate:"omitempty,max=50"`
}

// TechCount adalah hasil aggregation GET /portfolios/stats/technologies.
type TechCount struct {
	Technology string `json:"technology"`
	Total      int    `json:"total"`
}
