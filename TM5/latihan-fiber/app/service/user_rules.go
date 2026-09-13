package service

import (
	"strings"

	"latihan-fiber/app/model"
)

func ValidateCreate(req model.CreateUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "wajib diisi"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "wajib diisi"
	}
	if strings.TrimSpace(req.Password) == "" {
		errs["password"] = "wajib diisi"
	}

	return errs
}

func ValidateReplace(req model.ReplaceUserRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Password) == "" {
		errs["password"] = "wajib diisi pada PUT"
	}

	return errs
}

func ApplyPatch(current model.User, req model.PatchUserRequest) (model.User, map[string]string) {
	errs := map[string]string{}

	if req.Username != nil {
		if strings.TrimSpace(*req.Username) == "" {
			errs["username"] = "tidak boleh kosong"
		} else {
			current.Username = *req.Username // Ubah UserName -> Username
		}
	}
	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			errs["email"] = "tidak boleh kosong"
		} else {
			current.Email = *req.Email
		}
	}
	if req.Password != nil {
		if strings.TrimSpace(*req.Password) == "" {
			errs["password"] = "tidak boleh kosong"
		} else {
			current.Password = *req.Password
		}
	}

	return current, errs
}

func IsEmptyPatch(req model.PatchUserRequest) bool {
	return req.Username == nil && req.Email == nil && req.Password == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
