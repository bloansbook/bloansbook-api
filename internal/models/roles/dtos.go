package roles

import (
	"time"

	"github.com/google/uuid"
)

type RoleSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type PermissionSummary struct {
	ID     uuid.UUID `json:"id"`
	Code   string    `json:"code"`
	Module string    `json:"module"`
}

type CreateRolePayload struct {
	Name        string  `json:"name" validate:"required"`
	IsSystem    bool    `json:"isSystem" validate:"required"`
	Description *string `json:"description"`
}

type CreateRoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsSystem    bool      `json:"isSystem"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreatePermissionPayload struct {
	Code        string  `json:"code" validate:"required"`
	Module      string  `json:"module" validate:"required"`
	Description *string `json:"description"`
}

type CreatePermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Module      string    `json:"module"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateRolePermissionPayload struct {
	RoleID       uuid.UUID `json:"roleId" validate:"required"`
	PermissionID uuid.UUID `json:"permissionId" validate:"required"`
}

type CreateRolePermissionResponse struct {
	Role       RoleSummary       `json:"role"`
	Permission PermissionSummary `json:"permission"`
	CreatedAt  time.Time         `json:"createdAt"`
}

type RolePermissionsDTO struct {
	Role       RoleSummary         `json:"role"`
	Permission []PermissionSummary `json:"permission"`
}
