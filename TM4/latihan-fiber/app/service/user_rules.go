package service

import (
	"strings"

	"api-students/app/model"
)

func ValidateCreate(req model.CreateUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "wajib diisi"
	}

	return errs
}

func ValidateReplace(req model.ReplaceUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "wajib diisi pada PUT"
	}

	return errs
}

func ApplyPatch(current model.User, req model.PatchUserRequest) (model.User, map[string]string) {
	errs := map[string]string{}

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			current.Name = *req.Name
		}
	}
	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			errs["email"] = "tidak boleh kosong"
		} else {
			current.Email = *req.Email
		}
	}

	return current, errs
}

func IsEmptyPatch(req model.PatchUserRequest) bool {
	return req.Name == nil && req.Email == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
