package model

type Supervisor struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Quota           int    `json:"quota"`
	SupervisedCount int    `json:"supervised_count"`
	Version         int    `json:"version"`
}

type AssignSupervisionRequest struct {
	StudentID int `json:"student_id" validate:"required,gt=0"`
}

type SwitchSupervisionRequest struct {
	StudentID       int `json:"student_id" validate:"required,gt=0"`
	NewSupervisorID int `json:"new_supervisor_id" validate:"required,gt=0"`
}

type UpdateSupervisorQuotaRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitnil,notblank,max=100"`
	Quota           *int    `json:"quota,omitempty" validate:"omitnil,gte=0"`
	ExpectedVersion int     `json:"expected_version" validate:"required,gte=1"`
}
