package suppliers

import (
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

// SupplierDTO is the full supplier response shape returned to API consumers.
type SupplierDTO struct {
	ID uuid.UUID `json:"id"`

	SupplierID string `json:"supplierId"`

	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Email   *string `json:"email"`
	Address *string `json:"address"`

	Notes    *string          `json:"notes"`
	Category SupplierCategory `json:"category"`
	Currency string           `json:"currency"`

	Status models.BaseStatus `json:"status"`

	CreatedBy staff.StaffSummary `json:"createdBy"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SupplierSummary is the compact supplier shape used inside responses.
type SupplierSummary struct {
	SupplierID string            `json:"supplierId"`
	Name       string            `json:"name"`
	Phone      string            `json:"phone"`
	Email      *string           `json:"email,omitempty"`
	Category   SupplierCategory  `json:"category"`
	Status     models.BaseStatus `json:"status"`
}

// CreateSupplierPayload is the request body for creating a supplier.
// currency is intentionally omitted — v1 is NGN-only (DB default + CHECK).
type CreateSupplierPayload struct {
	Name     string           `json:"name"     validate:"required"`
	Phone    string           `json:"phone"    validate:"required"`
	Email    *string          `json:"email"`
	Address  *string          `json:"address"`
	Notes    *string          `json:"notes"`
	Category SupplierCategory `json:"category" validate:"required,oneof=raw_materials printing logistics artisans utilities other"`
}

// UpdateSupplierPayload is the request body for partial supplier updates.
// supplierId, currency and createdBy are immutable and cannot be changed here.
type UpdateSupplierPayload struct {
	Name     *string            `json:"name,omitempty"`
	Phone    *string            `json:"phone,omitempty"`
	Email    *string            `json:"email,omitempty"`
	Address  *string            `json:"address,omitempty"`
	Notes    *string            `json:"notes,omitempty"`
	Category *SupplierCategory  `json:"category,omitempty"`
	Status   *models.BaseStatus `json:"status,omitempty"`
}

// CreateSupplierResponse is returned after a successful supplier creation.
type CreateSupplierResponse struct {
	ID        uuid.UUID       `json:"id"`
	Supplier  SupplierSummary `json:"supplier"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// UpdateSupplierResponse is returned after a successful supplier update.
type UpdateSupplierResponse struct {
	ID        uuid.UUID       `json:"id"`
	Supplier  SupplierSummary `json:"supplier"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// SupplierFilter holds all optional query parameters for the list suppliers
// endpoint. Zero values mean "no filter applied" for that field.
type SupplierFilter struct {
	Search    string // matches name, supplier_id, phone (ILIKE)
	Category  string // one of: raw_materials, printing, logistics, artisans, utilities, other
	Status    string // one of: active, inactive
	SortBy    string // createdAt, name, supplierId, category, status
	SortOrder string // asc or desc (default: desc)
	Limit     int
	Offset    int
}
