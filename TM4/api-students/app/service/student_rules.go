package service

import (
	"strings"

	"api-students/app/model"
)

// Validasi untuk POST (Create)
func ValidateCreateStudent(req model.CreateStudentRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "NIM wajib diisi"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "nama wajib diisi"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "email wajib diisi"
	}

	return errs
}

// Validasi untuk PUT (Replace)
func ValidateReplaceStudent(req model.ReplaceStudentRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "NIM wajib diisi"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "nama wajib diisi"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "email wajib diisi"
	}

	return errs
}

// Penerapan perubahan untuk PATCH (Partial Update)
func ApplyPatchStudent(current model.Student, req model.PatchStudentRequest) (model.Student, map[string]string) {
	errs := make(map[string]string)
	updated := current

	if req.NIM != nil {
		val := strings.TrimSpace(*req.NIM)
		if val == "" {
			errs["nim"] = "NIM tidak boleh kosong"
		} else {
			updated.NIM = val
		}
	}

	if req.Name != nil {
		val := strings.TrimSpace(*req.Name)
		if val == "" {
			errs["name"] = "nama tidak boleh kosong"
		} else {
			updated.Name = val
		}
	}

	if req.Email != nil {
		val := strings.TrimSpace(*req.Email)
		if val == "" {
			errs["email"] = "email tidak boleh kosong"
		} else {
			updated.Email = val
		}
	}

	return updated, errs
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 || total <= 0 {
		return 0
	}
	pages := total / limit
	if total%limit != 0 {
		pages++
	}
	return pages
}
