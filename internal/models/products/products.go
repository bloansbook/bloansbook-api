package products

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	models.BaseModel

	ProductID string

	Name        string
	Description *string

	Status models.BaseStatus

	CreatedBy uuid.UUID
}

type ProductVariants struct {
	models.BaseModel

	ProductID uuid.UUID

	SKU string

	Size  *string
	Color *string

	Attributes []byte

	SellingPrice decimal.Decimal

	CurrentQuantity   int
	CurrentWAC        decimal.Decimal
	CurrentStockValue decimal.Decimal

	Status models.BaseStatus

	CreatedBy uuid.UUID
}
