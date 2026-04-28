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
