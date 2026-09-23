package customers_test

import (
	"testing"
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/customers"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

func sampleCustomer() customers.Customers {
	now := time.Date(2026, 8, 1, 9, 0, 0, 0, time.UTC)
	email := "ada@example.com"
	creatorEmail := "funmi@example.com"

	return customers.Customers{
		BaseModel: models.BaseModel{
			BaseWithId:        models.BaseWithId{ID: uuid.New()},
			BaseWithCreatedAt: models.BaseWithCreatedAt{CreatedAt: now},
			BaseWithUpdatedAt: models.BaseWithUpdatedAt{UpdatedAt: now},
		},
		CustomerID:        "CUS-0001",
		Name:              "Ada Lovelace",
		Phone:             "+2348012345678",
		Email:             &email,
		Type:              customers.CustomerTypeRetail,
		Currency:          "NGN",
		Status:            models.StatusActive,
		CreatorID:         uuid.New(),
		CreatorStaffID:    "BLN-0002",
		CreatorFirstName:  "Funmi",
		CreatorLastName:   "Adeyemi",
		CreatorEmail:      &creatorEmail,
		CreatorDepartment: "sales",
		CreatorJobTitle:   "Sales Executive",
		CreatorStatus:     staff.StaffStatusActive,
	}
}

func TestCustomers_ToSummary(t *testing.T) {
	c := sampleCustomer()
	s := c.ToSummary()

	if s.CustomerID != c.CustomerID {
		t.Errorf("CustomerID: got %q want %q", s.CustomerID, c.CustomerID)
	}
	if s.Name != c.Name {
		t.Errorf("Name: got %q want %q", s.Name, c.Name)
	}
	if s.Phone != c.Phone {
		t.Errorf("Phone: got %q want %q", s.Phone, c.Phone)
	}
	if s.Type != c.Type {
		t.Errorf("Type: got %q want %q", s.Type, c.Type)
	}
	if s.Status != c.Status {
		t.Errorf("Status: got %q want %q", s.Status, c.Status)
	}
	if s.Email == nil || *s.Email != *c.Email {
		t.Errorf("Email not mapped: got %v", s.Email)
	}
}

func TestCustomers_ToDTO_EmbedsCreatorSummary(t *testing.T) {
	c := sampleCustomer()
	dto := c.ToDTO()

	if dto.ID != c.ID {
		t.Errorf("ID: got %v want %v", dto.ID, c.ID)
	}
	if dto.CustomerID != c.CustomerID {
		t.Errorf("CustomerID: got %q want %q", dto.CustomerID, c.CustomerID)
	}
	if dto.Currency != "NGN" {
		t.Errorf("Currency: got %q want NGN", dto.Currency)
	}
	if dto.Status != c.Status {
		t.Errorf("Status: got %q want %q", dto.Status, c.Status)
	}

	// The creator is embedded as a staff summary sourced from the Creator* columns.
	if dto.CreatedBy.StaffID != c.CreatorStaffID {
		t.Errorf("CreatedBy.StaffID: got %q want %q", dto.CreatedBy.StaffID, c.CreatorStaffID)
	}
	if dto.CreatedBy.FirstName != c.CreatorFirstName {
		t.Errorf("CreatedBy.FirstName: got %q want %q", dto.CreatedBy.FirstName, c.CreatorFirstName)
	}
	if dto.CreatedBy.Status != c.CreatorStatus {
		t.Errorf("CreatedBy.Status: got %q want %q", dto.CreatedBy.Status, c.CreatorStatus)
	}
	if dto.CreatedAt != c.CreatedAt || dto.UpdatedAt != c.UpdatedAt {
		t.Errorf("timestamps not carried through: created=%v updated=%v", dto.CreatedAt, dto.UpdatedAt)
	}
}
