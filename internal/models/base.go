package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseWithId struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type BaseWithCreatedAt struct {
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type BaseWithUpdatedAt struct {
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type BaseModel struct {
	BaseWithId
	BaseWithCreatedAt
	BaseWithUpdatedAt
}

type BaseModelWithoutUpdatedAt struct {
	BaseWithId
	BaseWithCreatedAt
}

type BaseStatus string

const (
	StatusActive   BaseStatus = "active"
	StatusInactive BaseStatus = "inactive"
)

type DataWithPagination struct {
	Data       interface{} `json:"data_items"`
	Count      int         `json:"count"`
	TotalCount int         `json:"total_count"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
}
