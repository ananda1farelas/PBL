package model

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

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	SortBy   string
	Order    string
	IsActive *bool
}

func (q ListQuery) Offset() int {
	if q.Page <= 1 {
		return 0
	}
	return (q.Page - 1) * q.Limit
}
