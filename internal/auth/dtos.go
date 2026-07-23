package auth

import (
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

type ProfileDTO struct {
	Staff Staff  `json:"me"`
	Roles []Role `json:"roles"`
}

// TODO: Setup Login with staffID and Password with supabase auth
// TODO: Return Response which includes the token - should save the role as part of the payload
// TODO: Check the payload to see know role and role_permission
// TODO: Authorise based on the role and permissions

// TODO: Setup resend to send email to newly added users, their password and staffID so that they can login
