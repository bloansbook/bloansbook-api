package customers

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

// CustomerDTO is the full customer response shape returned to API consumers.
type CustomerDTO struct {
	ID uuid.UUID `json:"id"`

	CustomerID string `json:"customerId"`

	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`

	Notes    *string      `json:"notes"`
	Type     CustomerType `json:"type"`
	Currency string       `json:"currency"`

	Status models.BaseStatus `json:"status"`

	CreatedBy staff.StaffSummary `json:"createdBy"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CustomerSummary is the compact customer shape used inside responses.
type CustomerSummary struct {
	CustomerID string            `json:"customerId"`
	Name       string            `json:"name"`
	Phone      string            `json:"phone"`
	Email      *string           `json:"email,omitempty"`
	Type       CustomerType      `json:"type"`
	Status     models.BaseStatus `json:"status"`
}

// CreateCustomerPayload is the request body for creating a customer.
// currency is intentionally omitted — v1 is NGN-only (DB default + CHECK).
type CreateCustomerPayload struct {
	Name    string       `json:"name"    validate:"required"`
	Phone   string       `json:"phone"   validate:"required"`
	Email   *string      `json:"email"`
	Address *string      `json:"address"`
	Notes   *string      `json:"notes"`
	Type    CustomerType `json:"type"    validate:"required,oneof=retail corporate"`
}

// UpdateCustomerPayload is the request body for partial customer updates.
// customerId, currency and createdBy are immutable and cannot be changed here.
type UpdateCustomerPayload struct {
	Name    *string            `json:"name,omitempty"`
	Phone   *string            `json:"phone,omitempty"`
	Email   *string            `json:"email,omitempty"`
	Address *string            `json:"address,omitempty"`
	Notes   *string            `json:"notes,omitempty"`
	Type    *CustomerType      `json:"type,omitempty"`
	Status  *models.BaseStatus `json:"status,omitempty"`
}

// CreateCustomerResponse is returned after a successful customer creation.
type CreateCustomerResponse struct {
	ID        uuid.UUID       `json:"id"`
	Customer  CustomerSummary `json:"customer"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// UpdateCustomerResponse is returned after a successful customer update.
type UpdateCustomerResponse struct {
	ID        uuid.UUID       `json:"id"`
	Customer  CustomerSummary `json:"customer"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// CustomerFilter holds all optional query parameters for the list customers
// endpoint. Zero values mean "no filter applied" for that field.
type CustomerFilter struct {
	Search    string // matches name, customer_id, phone (ILIKE)
	Type      string // one of: retail, corporate
	Status    string // one of: active, inactive
	SortBy    string // createdAt, name, customerId, type, status
	SortOrder string // asc or desc (default: desc)
	Limit     int
	Offset    int
}
