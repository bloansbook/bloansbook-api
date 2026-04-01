package customers

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/google/uuid"
)

type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "retail"
	CustomerTypeCompany    CustomerType = "corporate"
)

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

	CreatedBy uuid.UUID `json:"createdBy" db:"created_by"`

	Status models.BaseStatus `json:"status" db:"status"`
}
