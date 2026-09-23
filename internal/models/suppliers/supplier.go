package suppliers

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

type SupplierCategory string

const (
	SupplierRawMaterials SupplierCategory = "raw_materials"
	SupplierPrinting     SupplierCategory = "printing"
	SupplierLogistics    SupplierCategory = "logistics"
	SupplierArtisans     SupplierCategory = "artisans"
	SupplierUtilities    SupplierCategory = "utilities"
	SupplierOther        SupplierCategory = "other"
)

// Supplier is the full supplier entity. The Creator* fields are populated via a
// LEFT JOIN on staff(created_by) so responses can embed a creator summary.
type Supplier struct {
	models.BaseModel

	SupplierID string `json:"supplierId" db:"supplier_id"`

	Name    string  `json:"name" db:"name"`
	Phone   string  `json:"phone" db:"phone"`
	Email   *string `json:"email" db:"email"`
	Address *string `json:"address" db:"address"`

	Notes    *string          `json:"notes" db:"notes"`
	Category SupplierCategory `json:"category" db:"category"`
	Currency string           `json:"currency" db:"currency"`

	Status models.BaseStatus `json:"status" db:"status"`

	CreatorID         uuid.UUID         `db:"creator_id"`
	CreatorStaffID    string            `db:"creator_staff_id"`
	CreatorFirstName  string            `db:"creator_first_name"`
	CreatorLastName   string            `db:"creator_last_name"`
	CreatorEmail      *string           `db:"creator_email"`
	CreatorPhone      *string           `db:"creator_phone"`
	CreatorDepartment string            `db:"creator_department"`
	CreatorJobTitle   string            `db:"creator_job_title"`
	CreatorStatus     staff.StaffStatus `db:"creator_status"`
}

// ToSummary converts a Supplier to its compact summary shape.
func (s *Supplier) ToSummary() SupplierSummary {
	return SupplierSummary{
		SupplierID: s.SupplierID,
		Name:       s.Name,
		Phone:      s.Phone,
		Email:      s.Email,
		Category:   s.Category,
		Status:     s.Status,
	}
}

// ToDTO converts a Supplier to the full SupplierDTO response shape,
// embedding the creator as a staff summary.
func (s *Supplier) ToDTO() SupplierDTO {
	return SupplierDTO{
		ID:         s.ID,
		SupplierID: s.SupplierID,
		Name:       s.Name,
		Phone:      s.Phone,
		Email:      s.Email,
		Address:    s.Address,
		Notes:      s.Notes,
		Category:   s.Category,
		Currency:   s.Currency,
		Status:     s.Status,
		CreatedBy: staff.StaffSummary{
			StaffID:    s.CreatorStaffID,
			FirstName:  s.CreatorFirstName,
			LastName:   s.CreatorLastName,
			Email:      s.CreatorEmail,
			Phone:      s.CreatorPhone,
			Department: s.CreatorDepartment,
			JobTitle:   s.CreatorJobTitle,
			Status:     s.CreatorStatus,
		},
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
