package model

type Nilai struct {
	ID             int     `json:"id"`
	IDStudent      int     `json:"id_student"`
	NamaMataKuliah string  `json:"nama_mata_kuliah"`
	Nilai          float64 `json:"nilai"`
}

type CreateNilaiRequest struct {
	IDStudent      int     `json:"id_student"`
	NamaMataKuliah string  `json:"nama_mata_kuliah"`
	Nilai          float64 `json:"nilai"`
}

type UpdateNilaiRequest struct {
	IDStudent      int     `json:"id_student"`
	NamaMataKuliah string  `json:"nama_mata_kuliah"`
	Nilai          float64 `json:"nilai"`
}
