package service

import (
	"strings"

	"api-students/app/model"
)

// ApplyPatchStudent menerapkan perubahan dari request PATCH.
// Validasi request dilakukan menggunakan helper.ValidateStruct.
func ApplyPatchStudent(current model.Student, req model.PatchStudentRequest) model.Student {
	updated := current

	if req.NIM != nil {
		updated.NIM = strings.TrimSpace(*req.NIM)
	}

	if req.Name != nil {
		updated.Name = strings.TrimSpace(*req.Name)
	}

	if req.Grade != nil {
		updated.Grade = strings.TrimSpace(*req.Grade)
	}

	if req.IsActive != nil {
		updated.IsActive = *req.IsActive
	}

	return updated
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
