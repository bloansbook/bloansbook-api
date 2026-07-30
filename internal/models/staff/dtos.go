package staff

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StaffDTO is the full staff response shape returned to API consumers.
type StaffDTO struct {
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

	Status  StaffStatus `json:"status"`
	FiredAt *time.Time  `json:"firedAt,omitempty"`

	HasLogin bool `json:"hasLogin"`

	CreatedBy StaffSummary        `json:"createdBy"`
	Roles     []roles.RoleSummary `json:"roles"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StaffSummary is the compact staff shape used inside responses.
type StaffSummary struct {
	StaffID    string      `json:"staffId"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Email      *string     `json:"email,omitempty"`
	Phone      *string     `json:"phone,omitempty"`
	Department string      `json:"department"`
	JobTitle   string      `json:"jobTitle"`
	Status     StaffStatus `json:"status"`
}

// CreateStaffPayload is the request body for creating a staff member.
type CreateStaffPayload struct {
	FirstName   string     `json:"firstName"   validate:"required"`
	LastName    string     `json:"lastName"    validate:"required"`
	Email       *string    `json:"email"`
	Phone       *string    `json:"phone"`
	Address     *string    `json:"address"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	DateOfHire  time.Time  `json:"dateOfHire"  validate:"required"`

	EmergencyContactName  *string `json:"emergencyContactName"`
	EmergencyContactPhone *string `json:"emergencyContactPhone"`

	BankName          *string `json:"bankName"`
	BankAccountNumber *string `json:"bankAccountNumber"`
	BankAccountName   *string `json:"bankAccountName"`

	Department string          `json:"department" validate:"required"`
	JobTitle   string          `json:"jobTitle"   validate:"required"`
	PayType    string          `json:"payType"    validate:"required"`
	BaseSalary decimal.Decimal `json:"baseSalary" validate:"required"`

	HasLogin    bool    `json:"hasLogin"`
	SupabaseUID *string `json:"supabaseUID"`

	Status StaffStatus `json:"status" validate:"required,oneof=active inactive fired"`
}

// Credentials holds the generated login credentials before hashing/storage.
type Credentials struct {
	StaffID  string `json:"staffId"   validate:"required"`
	Password string `json:"password"  validate:"required"`
}

// CreateStaffResponse is returned after a successful staff creation.
type CreateStaffResponse struct {
	ID        uuid.UUID    `json:"id"`
	Staff     StaffSummary `json:"staff_info"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// UpdateStaffPayload is the request body for partial staff updates.
type UpdateStaffPayload struct {
	FirstName             *string          `json:"firstName,omitempty"`
	LastName              *string          `json:"lastName,omitempty"`
	Email                 *string          `json:"email,omitempty"`
	Phone                 *string          `json:"phone,omitempty"`
	Address               *string          `json:"address,omitempty"`
	DateOfBirth           *time.Time       `json:"dateOfBirth,omitempty"`
	EmergencyContactName  *string          `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone *string          `json:"emergencyContactPhone,omitempty"`
	BankName              *string          `json:"bankName,omitempty"`
	BankAccountNumber     *string          `json:"bankAccountNumber,omitempty"`
	BankAccountName       *string          `json:"bankAccountName,omitempty"`
	Department            *string          `json:"department,omitempty"`
	JobTitle              *string          `json:"jobTitle,omitempty"`
	PayType               *string          `json:"payType,omitempty"`
	BaseSalary            *decimal.Decimal `json:"baseSalary,omitempty"`
	Status                *StaffStatus     `json:"status,omitempty"`
}

// UpdateStaffResponse is returned after a successful staff update.
type UpdateStaffResponse struct {
	ID        uuid.UUID    `json:"id"`
	Staff     StaffSummary `json:"staff"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

// --- Staff Role DTOs ---

// AssignRolePayload is the request body for assigning a role to a staff member.
type AssignRolePayload struct {
	RoleID uuid.UUID `json:"roleId" validate:"required"`
	Reason *string   `json:"reason"`
}

// RevokeRolePayload is the request body for revoking a role from a staff member.
type RevokeRolePayload struct {
	RoleID uuid.UUID `json:"roleId" validate:"required"`
	Reason *string   `json:"reason"`
}

// UpdateRolePayload swaps the current role for a new one in a single transaction.
type UpdateRolePayload struct {
	OldRoleID uuid.UUID `json:"oldRoleId" validate:"required"`
	NewRoleID uuid.UUID `json:"newRoleId" validate:"required"`
	Reason    *string   `json:"reason"`
}

// StaffRoleResponse is returned after assign / revoke / update.
type StaffRoleResponse struct {
	StaffID  uuid.UUID `json:"staffId"`
	RoleID   uuid.UUID `json:"roleId"`
	RoleName string    `json:"roleName"`
	Action   string    `json:"action"`
}

// StaffRoleHistoryEntry is one row from staff_role_history.
type StaffRoleHistoryEntry struct {
	ID          uuid.UUID `json:"id"`
	RoleID      uuid.UUID `json:"roleId"`
	RoleName    string    `json:"roleName"`
	Action      string    `json:"action"`
	PerformedBy uuid.UUID `json:"performedBy"`
	Reason      *string   `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// StaffFilter holds all optional query parameters for the list staff endpoint.
// Zero values mean "no filter applied" for that field.
type StaffFilter struct {
	Search     string // matches first_name, last_name, staff_id (ILIKE)
	Status     string // one of: active, inactive, fired
	Department string // one of: factory, admin, sales, management
	SortBy     string // column to sort by: createdAt, firstName, lastName, staffId, department, status
	SortOrder  string // asc or desc (default: desc)
	Limit      int
	Offset     int
}

// --- Termination DTOs ---

// FireStaffPayload is the request body for terminating a staff member.
type FireStaffPayload struct {
	TerminationReason string `json:"terminationReason" validate:"required"`
}

// FireStaffResponse is returned after a successful staff termination.
type FireStaffResponse struct {
	ID                uuid.UUID `json:"id"`
	StaffID           uuid.UUID `json:"staffId"`
	TerminationReason string    `json:"terminationReason"`
	RecordedBy        uuid.UUID `json:"recordedBy"`
	RecordedAt        time.Time `json:"recordedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

// OverrideTerminationPayload is the request body for overriding a termination record.
// Only super_admin (staff.terminate permission) may call this endpoint.
type OverrideTerminationPayload struct {
	OverrideReason string `json:"overrideReason" validate:"required"`
}

// OverrideTerminationResponse is returned after a successful termination override.
type OverrideTerminationResponse struct {
	ID                uuid.UUID  `json:"id"`
	StaffID           uuid.UUID  `json:"staffId"`
	TerminationReason string     `json:"terminationReason"`
	IsOverridden      bool       `json:"isOverridden"`
	OverriddenBy      *uuid.UUID `json:"overriddenBy"`
	OverriddenAt      *time.Time `json:"overriddenAt"`
	OverrideReason    *string    `json:"overrideReason"`
	RecordedAt        time.Time  `json:"recordedAt"`
	CreatedAt         time.Time  `json:"createdAt"`
}
