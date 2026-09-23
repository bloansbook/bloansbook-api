package customers

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

type CustomerType string

const (
	CustomerTypeRetail    CustomerType = "retail"
	CustomerTypeCorporate CustomerType = "corporate"
)

// Customers is the full customer entity. The Creator* fields are populated via a
// LEFT JOIN on staff(created_by) so responses can embed a creator summary.
type Customers struct {
	models.BaseModel

	CustomerID string `json:"customerId" db:"customer_id"`

	Name    string  `json:"name" db:"name"`
	Phone   string  `json:"phone" db:"phone"`
	Email   *string `json:"email" db:"email"`
	Address *string `json:"address" db:"address"`

	Notes    *string      `json:"notes" db:"notes"`
	Type     CustomerType `json:"type" db:"type"`
	Currency string       `json:"currency" db:"currency"`

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

// ToSummary converts a Customers to its compact summary shape.
func (c *Customers) ToSummary() CustomerSummary {
	return CustomerSummary{
		CustomerID: c.CustomerID,
		Name:       c.Name,
		Phone:      c.Phone,
		Email:      c.Email,
		Type:       c.Type,
		Status:     c.Status,
	}
}

// ToDTO converts a Customers to the full CustomerDTO response shape,
// embedding the creator as a staff summary.
func (c *Customers) ToDTO() CustomerDTO {
	return CustomerDTO{
		ID:         c.ID,
		CustomerID: c.CustomerID,
		Name:       c.Name,
		Phone:      c.Phone,
		Email:      c.Email,
		Address:    c.Address,
		Notes:      c.Notes,
		Type:       c.Type,
		Currency:   c.Currency,
		Status:     c.Status,
		CreatedBy: staff.StaffSummary{
			StaffID:    c.CreatorStaffID,
			FirstName:  c.CreatorFirstName,
			LastName:   c.CreatorLastName,
			Email:      c.CreatorEmail,
			Phone:      c.CreatorPhone,
			Department: c.CreatorDepartment,
			JobTitle:   c.CreatorJobTitle,
			Status:     c.CreatorStatus,
		},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
