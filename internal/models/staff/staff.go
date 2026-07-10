package staff

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StaffStatus string
type RoleAction string

const (
	StaffStatusActive   StaffStatus = "active"
	StaffStatusInactive StaffStatus = "inactive"
	StaffStatusFired    StaffStatus = "fired"
)

const (
	RoleActionAssigned RoleAction = "assigned"
	RoleActionRevoked  RoleAction = "revoked"
)

type Staff struct {
	models.BaseModel

	StaffID  string  `json:"staffId" db:"staff_id"`
	Password *string `json:"_" db:"password_hash"`

	FirstName             string     `json:"firstName" db:"first_name"`
	LastName              string     `json:"lastName" db:"last_name"`
	Email                 *string    `json:"email" db:"email"`
	Phone                 *string    `json:"phone" db:"phone"`
	Address               *string    `json:"address" db:"address"`
	DateOfBirth           *time.Time `json:"dateOfBirth" db:"date_of_birth"`
	DateOfHire            time.Time  `json:"dateOfHire" db:"date_of_hire"`
	EmergencyContactName  *string    `json:"emergencyContactName" db:"emergency_contact_name"`
	EmergencyContactPhone *string    `json:"emergencyContactPhone" db:"emergency_contact_phone"`

	BankName          *string `json:"bankName" db:"bank_name"`
	BankAccountNumber *string `json:"bankAccountNumber" db:"bank_account_number"`
	BankAccountName   *string `json:"bankAccountName" db:"bank_account_name"`

	Department string          `json:"department" db:"department"`
	JobTitle   string          `json:"jobTitle" db:"job_title"`
	PayType    string          `json:"payType" db:"pay_type"`
	BaseSalary decimal.Decimal `json:"baseSalary" db:"base_salary"`

	Status  StaffStatus `json:"status" db:"status"`
	FiredAt *time.Time  `json:"firedAt,omitempty" db:"fired_at"`

	HasLogin    bool    `json:"hasLogin" db:"has_login"`
	SupabaseUID *string `json:"supabaseUID" db:"supabase_uid"`

	CreatorID         uuid.UUID   `db:"creator_id"`
	CreatorStaffID    string      `db:"creator_staff_id"`
	CreatorFirstName  string      `db:"creator_first_name"`
	CreatorLastName   string      `db:"creator_last_name"`
	CreatorEmail      *string     `db:"creator_email"`
	CreatorPhone      *string     `db:"creator_phone"`
	CreatorDepartment string      `db:"creator_department"`
	CreatorJobTitle   string      `db:"creator_job_title"`
	CreatorStatus     StaffStatus `db:"creator_status"`
}

type StaffCreate struct {
	models.BaseModel

	StaffID    string      `json:"staffId" db:"staff_id"`
	FirstName  string      `json:"firstName" db:"first_name"`
	LastName   string      `json:"lastName" db:"last_name"`
	Email      *string     `json:"email" db:"email"`
	Phone      *string     `json:"phone" db:"phone"`
	Department string      `json:"department" db:"department"`
	JobTitle   string      `json:"jobTitle" db:"job_title"`
	Status     StaffStatus `json:"status" db:"status"`
}

type FiredStaff struct {
	models.BaseModelWithoutUpdatedAt

	StaffID uuid.UUID `json:"staffId" db:"staff_id"`

	TerminationReason string    `json:"terminationReason" db:"termination_reason"`
	RecordedBy        uuid.UUID `json:"recordedBy" db:"recorded_by"`
	RecordedAt        time.Time `json:"recordedAt" db:"recorded_at"`

	IsOverridden   bool       `json:"isOverridden" db:"is_overridden"`
	OverriddenBy   *uuid.UUID `json:"overriddenBy" db:"overridden_by"`
	OverriddenAt   *time.Time `json:"overriddenAt" db:"overridden_at"`
	OverrideReason *string    `json:"overrideReason" db:"override_reason"`
}

type StaffRoles struct {
	StaffID uuid.UUID `json:"staffId" db:"staff_id"`
	RoleID  uuid.UUID `json:"roleId" db:"role_id"`

	AssignedAt time.Time `json:"assignedAt" db:"assigned_at"`
	AssignedBy uuid.UUID `json:"assignedBy" db:"assigned_by"`
}

type StaffRoleHistory struct {
	models.BaseModelWithoutUpdatedAt

	StaffID uuid.UUID `json:"staffId" db:"staff_id"`
	RoleID  uuid.UUID `json:"roleId" db:"role_id"`

	Action RoleAction `json:"action" db:"action"`

	PerformedBy uuid.UUID `json:"performedBy" db:"performed_by"`
	Reason      *string   `json:"reason,omitempty" db:"reason"`
}

// ToSummary converts a Staff to its compact summary shape.
func (s *Staff) ToSummary() StaffSummary {
	return StaffSummary{
		StaffID:    s.StaffID,
		FirstName:  s.FirstName,
		LastName:   s.LastName,
		Email:      s.Email,
		Phone:      s.Phone,
		Department: s.Department,
		JobTitle:   s.JobTitle,
		Status:     s.Status,
	}
}

// ToDTO converts a Staff to the full StaffDTO response shape.
// staffRoles should be passed in from a separate query; pass nil for an empty list.
func (s *Staff) ToDTO(staffRoles []roles.RoleSummary) StaffDTO {
	if staffRoles == nil {
		staffRoles = []roles.RoleSummary{}
	}
	return StaffDTO{
		ID:                    s.ID,
		StaffID:               s.StaffID,
		FirstName:             s.FirstName,
		LastName:              s.LastName,
		Email:                 s.Email,
		Phone:                 s.Phone,
		Address:               s.Address,
		DateOfBirth:           s.DateOfBirth,
		DateOfHire:            s.DateOfHire,
		EmergencyContactName:  s.EmergencyContactName,
		EmergencyContactPhone: s.EmergencyContactPhone,
		BankName:              s.BankName,
		BankAccountNumber:     s.BankAccountNumber,
		BankAccountName:       s.BankAccountName,
		Department:            s.Department,
		JobTitle:              s.JobTitle,
		PayType:               s.PayType,
		BaseSalary:            s.BaseSalary,
		Status:                s.Status,
		FiredAt:               s.FiredAt,
		HasLogin:              s.HasLogin,
		CreatedBy: StaffSummary{
			StaffID:    s.CreatorStaffID,
			FirstName:  s.CreatorFirstName,
			LastName:   s.CreatorLastName,
			Email:      s.CreatorEmail,
			Phone:      s.CreatorPhone,
			Department: s.CreatorDepartment,
			JobTitle:   s.CreatorJobTitle,
			Status:     s.CreatorStatus,
		},
		Roles:     staffRoles,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSummary converts a StaffCreate to its compact summary shape.
func (s *StaffCreate) ToSummary() StaffSummary {
	return StaffSummary{
		StaffID:    s.StaffID,
		FirstName:  s.FirstName,
		LastName:   s.LastName,
		Email:      s.Email,
		Phone:      s.Phone,
		Department: s.Department,
		JobTitle:   s.JobTitle,
		Status:     s.Status,
	}
}
