package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LoginDTO struct {
	StaffID  string `json:"staffId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Staff struct {
	ID uuid.UUID `json:"id"`

	StaffID string `json:"staffId"`

	FirstName             string     `json:"firstName"`
	LastName              string     `json:"lastName"`
	Email                 *string    `json:"email"`
	Phone                 *string    `json:"phone"`
	Address               *string    `json:"address"`
	DateOfBirth           *time.Time `json:"dateOfBirth"`
	DateOfHire            time.Time  `json:"dateOfHire"`
	EmergencyContactName  *string    `json:"emergencyContactName"`
	EmergencyContactPhone *string    `json:"emergencyContactPhone"`

	BankName          *string `json:"bankName"`
	BankAccountNumber *string `json:"bankAccountNumber"`
	BankAccountName   *string `json:"bankAccountName"`

	Department string          `json:"department"`
	JobTitle   string          `json:"jobTitle"`
	PayType    string          `json:"payType"`
	BaseSalary decimal.Decimal `json:"baseSalary"`

	Status  staff.StaffStatus `json:"status"`
	FiredAt *time.Time        `json:"firedAt,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ModulePermission struct {
	Module string   `json:"module"`
	Codes  []string `json:"codes"`
}

type Role struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Permissions []ModulePermission `json:"permissions"`
}

// RoleList is a JSON-scannable slice used when reading the roles column
// from a PostgreSQL json_agg result via pgx.
type RoleList []Role

func (r *RoleList) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for roles, got %T", src)
	}
	return json.Unmarshal(bytes, r)
}

// PermissionList is a JSON-scannable slice used when reading the flat
// permissions column from a PostgreSQL json_agg result via pgx.
type PermissionList []string

func (p *PermissionList) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for permissions, got %T", src)
	}
	return json.Unmarshal(bytes, p)
}

type ProfileDTO struct {
	Staff       Staff          `json:"me"`
	Roles       RoleList       `json:"roles"`
	Permissions PermissionList `json:"permissions"`
}
