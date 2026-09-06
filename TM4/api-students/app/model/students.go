package model

type Student struct {
	ID    int    `json:"id"`
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateStudentRequest struct {
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReplaceStudentRequest struct {
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PatchStudentRequest struct {
	NIM   *string `json:"nim,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
