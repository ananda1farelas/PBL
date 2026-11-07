package service

import (
	"errors"
	"fmt"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

type SupervisionService struct {
	tx           *repository.TxManager
	supervisors  repository.SupervisorRepository
	supervisions repository.SupervisionRepository
}

func NewSupervisionService(tx *repository.TxManager, supervisors repository.SupervisorRepository, supervisions repository.SupervisionRepository) *SupervisionService {
	return &SupervisionService{tx: tx, supervisors: supervisors, supervisions: supervisions}
}
func mapSupervisionError(err error) error {
	switch {
	case errors.Is(err, repository.ErrStudentNotFound), errors.Is(err, repository.ErrSupervisorNotFound), errors.Is(err, repository.ErrNotAssigned):
		return helper.NotFound(err.Error())
	case errors.Is(err, repository.ErrAlreadyAssigned), errors.Is(err, repository.ErrQuotaFull), errors.Is(err, repository.ErrVersionConflict), errors.Is(err, repository.ErrSameSupervisor):
		return helper.Conflict(err.Error())
	default:
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return helper.Conflict("mahasiswa sudah memiliki dosen pembimbing")
		}
		return fmt.Errorf("transaksi supervision gagal: %w", err)
	}
}
func (s *SupervisionService) Assign(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID dosen pembimbing tidak valid")
	}
	var req model.AssignSupervisionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}
	err := s.tx.WithTx(c.UserContext(), func(tx repository.DBTX) error {
		students := s.supervisions.WithTx(tx)
		supervisors := s.supervisors.WithTx(tx)
		if err := students.LockStudent(c.UserContext(), req.StudentID); err != nil {
			return err
		}
		if _, err := students.FindSupervisorIDByStudent(c.UserContext(), req.StudentID); err == nil {
			return repository.ErrAlreadyAssigned
		} else if !errors.Is(err, repository.ErrNotAssigned) {
			return err
		}
		if _, err := supervisors.FindByIDForUpdate(c.UserContext(), id); err != nil {
			return err
		}
		if err := supervisors.IncreaseCount(c.UserContext(), id); err != nil {
			return err
		}
		return students.Assign(c.UserContext(), req.StudentID, id)
	})
	if err != nil {
		return mapSupervisionError(err)
	}
	return helper.Created(c, "dosen pembimbing berhasil ditetapkan", fiber.Map{"student_id": req.StudentID, "supervisor_id": id})
}
func (s *SupervisionService) Switch(c *fiber.Ctx) error {
	var req model.SwitchSupervisionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}
	err := s.tx.WithTx(c.UserContext(), func(tx repository.DBTX) error {
		students := s.supervisions.WithTx(tx)
		supervisors := s.supervisors.WithTx(tx)
		if err := students.LockStudent(c.UserContext(), req.StudentID); err != nil {
			return err
		}
		oldID, err := students.FindSupervisorIDByStudent(c.UserContext(), req.StudentID)
		if err != nil {
			return err
		}
		if oldID == req.NewSupervisorID {
			return repository.ErrSameSupervisor
		}
		pair, err := supervisors.FindPairForUpdate(c.UserContext(), oldID, req.NewSupervisorID)
		if err != nil {
			return err
		}
		if pair[req.NewSupervisorID].SupervisedCount >= pair[req.NewSupervisorID].Quota {
			return repository.ErrQuotaFull
		}
		if err := supervisors.IncreaseCount(c.UserContext(), req.NewSupervisorID); err != nil {
			return err
		}
		if err := students.Switch(c.UserContext(), req.StudentID, req.NewSupervisorID, oldID); err != nil {
			return err
		}
		return supervisors.DecreaseCount(c.UserContext(), oldID)
	})
	if err != nil {
		return mapSupervisionError(err)
	}
	return helper.OK(c, "dosen pembimbing berhasil dipindahkan", fiber.Map{"student_id": req.StudentID, "supervisor_id": req.NewSupervisorID})
}
func (s *SupervisionService) UpdateSupervisor(c *fiber.Ctx) error {
	id, ok := helper.ParamID(c)
	if !ok {
		return helper.BadRequest("ID dosen pembimbing tidak valid")
	}
	var req model.UpdateSupervisorQuotaRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("format payload JSON tidak valid")
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = &trimmed
	}
	if errs := helper.ValidateStruct(req); len(errs) > 0 {
		return helper.ValidationFailed("validasi gagal", errs)
	}
	var updated model.Supervisor
	err := s.tx.WithTx(c.UserContext(), func(tx repository.DBTX) error {
		var err error
		updated, err = s.supervisors.WithTx(tx).UpdateOptimistic(c.UserContext(), id, req)
		return err
	})
	if err != nil {
		return mapSupervisionError(err)
	}
	return helper.OK(c, "data dosen pembimbing berhasil diperbarui", updated)
}
