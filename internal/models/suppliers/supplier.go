package suppliers

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
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

type Supplier struct {
	models.BaseModel

	SupplierID string

	Name    string
	Phone   string
	Email   *string
	Address *string

	Currency string
	Category SupplierCategory

	Notes *string

	CreatedBy uuid.UUID

	Status models.BaseStatus
}
