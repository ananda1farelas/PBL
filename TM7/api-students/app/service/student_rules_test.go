package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestCreateStudentValidation(t *testing.T) {
	t.Run("input valid", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:      "18221001",
			Name:     "Budi Santoso",
			Grade:    "85",
			IsActive: true,
		}

		errs := helper.ValidateStruct(req)
		if len(errs) != 0 {
			t.Errorf("diharapkan tidak ada error, dapat: %v", errs)
		}
	})

	t.Run("NIM tidak valid", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:   "ABC123",
			Name:  "Budi Santoso",
			Grade: "85",
		}

		errs := helper.ValidateStruct(req)
		if errs["nim"] == "" {
			t.Error("diharapkan ada error pada NIM")
		}
	})

	t.Run("nama kosong", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:   "18221001",
			Name:  "",
			Grade: "85",
		}

		errs := helper.ValidateStruct(req)
		if errs["name"] == "" {
			t.Error("diharapkan ada error pada nama")
		}
	})
}

func TestReplaceStudentValidation(t *testing.T) {
	req := model.ReplaceStudentRequest{
		NIM:   "18221001",
		Name:  "",
		Grade: "85",
	}

	errs := helper.ValidateStruct(req)
	if errs["name"] == "" {
		t.Error("diharapkan ada error pada nama")
	}
}

func TestPatchStudentValidation(t *testing.T) {
	t.Run("nama kosong ditolak", func(t *testing.T) {
		emptyName := ""
		req := model.PatchStudentRequest{
			Name: &emptyName,
		}

		errs := helper.ValidateStruct(req)
		if errs["name"] == "" {
			t.Error("nama kosong seharusnya ditolak")
		}
	})

	t.Run("nama tidak dikirim diterima", func(t *testing.T) {
		req := model.PatchStudentRequest{
			IsActive: boolPointer(false),
		}

		errs := helper.ValidateStruct(req)
		if len(errs) != 0 {
			t.Errorf("field yang tidak dikirim seharusnya dilewati: %v", errs)
		}
	})
}

func TestApplyPatchStudent(t *testing.T) {
	current := model.Student{
		ID:       1,
		NIM:      "18221001",
		Name:     "Budi Lama",
		Grade:    "80",
		IsActive: true,
	}

	newName := "Budi Baru"
	inactive := false

	req := model.PatchStudentRequest{
		Name:     &newName,
		IsActive: &inactive,
	}

	updated := ApplyPatchStudent(current, req)

	if updated.Name != "Budi Baru" {
		t.Errorf("nama tidak berubah: %s", updated.Name)
	}

	if updated.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}

	if updated.NIM != current.NIM {
		t.Error("NIM yang tidak dikirim seharusnya tidak berubah")
	}
}

func boolPointer(value bool) *bool {
	return &value
}
