package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateCreateStudent(t *testing.T) {
	t.Run("input valid", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:   "18221001",
			Name:  "Budi Santoso",
			Grade: "85",
		}
		errs := ValidateCreateStudent(req)
		if len(errs) != 0 {
			t.Errorf("diharapkan 0 error, dapat: %v", errs)
		}
	})

	t.Run("input kosong", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:   "",
			Name:  "   ",
		}
		errs := ValidateCreateStudent(req)
		if len(errs) != 3 {
			t.Errorf("diharapkan 3 error, dapat: %d", len(errs))
		}
	})
}

func TestValidateReplaceStudent(t *testing.T) {
	t.Run("input kosong sebagian", func(t *testing.T) {
		req := model.ReplaceStudentRequest{
			NIM:   "18221001",
			Name:  "",
			Grade: "85",
		}
		errs := ValidateReplaceStudent(req)
		if len(errs) != 1 || errs["name"] == "" {
			t.Errorf("diharapkan error pada nama, dapat: %v", errs)
		}
	})
}

func TestApplyPatchStudent(t *testing.T) {
	t.Run("patch dengan string kosong", func(t *testing.T) {
		current := model.Student{ID: 1, NIM: "18221001", Name: "Lama"}
		invalidName := "   "
		req := model.PatchStudentRequest{Name: &invalidName}

		_, errs := ApplyPatchStudent(current, req)
		if len(errs) != 1 || errs["name"] == "" {
			t.Errorf("diharapkan error validasi nama kosong, dapat: %v", errs)
		}
	})
}
