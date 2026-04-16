package roles

import (
	"encoding/json"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/google/uuid"
)

type PermissionSummaryList []PermissionSummary

func (p *PermissionSummaryList) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for permissions, got %T", src)
	}
	return json.Unmarshal(bytes, p)
}

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

	RoleID           uuid.UUID `json:"roleId" db:"role_id"`
	PermissionID     uuid.UUID `json:"permissionId" db:"permission_id"`
	Role             string    `json:"role" db:"role_name"`
	PermissionCode   string    `json:"permissionCode" db:"permission_code"`
	PermissionModule string    `json:"permissionModule" db:"permission_module"`
}

type RoleWithPermissions struct {
	ID          uuid.UUID             `json:"id" db:"id"`
	Name        string                `json:"name" db:"name"`
	Permissions PermissionSummaryList `json:"permissions" db:"permissions"`
}
