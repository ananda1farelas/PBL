package model

type WebResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    any               `json:"data,omitempty"`
	Meta    any               `json:"meta,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Meta menyimpan metadata pagination berbasis nomor halaman.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Cursor menyimpan posisi terakhir yang dikirim ke client.
// Pasangan created_at dan id membuat urutan tetap unik.
type Cursor struct {
	CreatedAt string
	ID        int
}

// CursorQuery adalah parameter pagination berbasis keyset.
type CursorQuery struct {
	Limit int
	After *Cursor
}

// CursorMeta adalah metadata pagination tanpa total halaman.
type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
