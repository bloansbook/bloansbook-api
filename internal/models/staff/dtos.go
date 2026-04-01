package staff

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StaffDTO struct {
	ID uuid.UUID `json:"id"`

	StaffID string `json:"staffId" db:"staff_id"`

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

	HasLogin bool `json:"hasLogin" db:"has_login"`

	CreatedBy StaffSummary `json:"createdBy" db:"created_by"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type StaffSummary struct {
	StaffID   string      `json:"staffId"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     *string     `json:"email,omitempty"`
	Phone     *string     `json:"phone,omitempty"`
	Status    StaffStatus `json:"status"`
}

type CreateStaffPayload struct {
	StaffID  string `json:"staffId" validate:"required"`
	Password string `json:"password" validate:"required"`

	FirstName   string     `json:"firstName" validate:"required"`
	LastName    string     `json:"lastName" validate:"required"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	Address     *string    `json:"address"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	DateOfHire  time.Time  `json:"dateOfHire" validate:"required"`

	EmergencyContactName  *string `json:"emergencyContactName"`
	EmergencyContactPhone *string `json:"emergencyContactPhone"`

	BankName          *string `json:"bankName"`
	BankAccountNumber *string `json:"bankAccountNumber"`
	BankAccountName   *string `json:"bankAccountName"`

	Department string          `json:"department" validate:"required"`
	JobTitle   string          `json:"jobTitle" validate:"required"`
	PayType    string          `json:"payType" validate:"required"`
	BaseSalary decimal.Decimal `json:"baseSalary" validate:"required"`

	HasLogin bool `json:"hasLogin"`

	Status StaffStatus `json:"status" validate:"required,oneof=active inactive fired"`
}

type CreateStaffResponse struct {
	ID    uuid.UUID    `json:"id"`
	Staff StaffSummary `json:"staff"`

	CreatedBy StaffSummary `json:"createdBy"`

	Roles       []roles.RoleSummary       `json:"roles"`
	Permissions []roles.PermissionSummary `json:"permissions"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
