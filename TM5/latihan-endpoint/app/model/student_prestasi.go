package model

type StudentPrestasi struct {
	ID           int    `json:"id"`
	IDPeserta    int    `json:"idpeserta"`
	IDPrestasi   int    `json:"idprestasi"`
	Juara        string `json:"juara"`
	NamaPrestasi string `json:"namaprestasi"`
}

type CreateStudentPrestasiRequest struct {
	IDPeserta    int    `json:"idpeserta"`
	IDPrestasi   int    `json:"idprestasi"`
	Juara        string `json:"juara"`
	NamaPrestasi string `json:"namaprestasi"`
}

type UpdateStudentPrestasiRequest struct {
	IDPeserta    int    `json:"idpeserta"`
	IDPrestasi   int    `json:"idprestasi"`
	Juara        string `json:"juara"`
	NamaPrestasi string `json:"namaprestasi"`
}

type PatchStudentPrestasiRequest struct {
	IDPeserta    *int    `json:"idpeserta,omitempty"`
	IDPrestasi   *int    `json:"idprestasi,omitempty"`
	Juara        *string `json:"juara,omitempty"`
	NamaPrestasi *string `json:"namaprestasi,omitempty"`
}
