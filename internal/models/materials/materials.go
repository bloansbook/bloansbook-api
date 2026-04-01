package materials

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Materials struct {
	models.BaseModel

	MaterialID string

	Name string

	UnitOfMeasure string

	ReorderLevel int

	CurrentQuantity   int
	CurrentWAC        decimal.Decimal
	CurrentStockValue decimal.Decimal

	Status models.BaseStatus

	CreatedBy uuid.UUID
}
