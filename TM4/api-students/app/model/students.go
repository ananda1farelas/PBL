package model

type Student struct {
	ID        int    `json:"id"`
	NIM       string `json:"nim"`
	Name      string `json:"name"`
	Grade     string `json:"grade"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
}

type CreateStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive bool   `json:"isActive"`
}

type ReplaceStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive bool   `json:"isActive"`
}

type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty"`
	Name     *string `json:"name,omitempty"`
	Grade    *string `json:"grade,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}
