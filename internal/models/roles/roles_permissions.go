package roles

import (
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/google/uuid"
)

type Roles struct {
	models.BaseModel

	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsSystem    bool   `json:"isSystem" db:"is_system"`
}

type Permissions struct {
	models.BaseModelWithoutUpdatedAt

	Code   string `json:"code" db:"code"`
	Module string `json:"module" db:"module"`

	Description string `json:"description" db:"description"`
}

type RolePermissions struct {
	models.BaseWithCreatedAt

	RoleID       uuid.UUID `json:"roleId" db:"role_id"`
	PermissionID uuid.UUID `json:"permissionId" db:"permission_id"`
}
