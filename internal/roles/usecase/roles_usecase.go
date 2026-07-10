package usecase

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/bloansbook/bloansbook-api/internal/roles/repository"
	"github.com/google/uuid"
)

type RolesUsecase struct {
	repository *repository.RolesRepository
}

func NewRolesUsecase(repo *repository.RolesRepository) *RolesUsecase {
	return &RolesUsecase{repository: repo}
}

func (u *RolesUsecase) CreateRole(ctx context.Context, payload *roles.CreateRolePayload) (*roles.Roles, error) {
	if payload.Name == "" {
		return nil, fmt.Errorf("role name is required")
	}
	return u.repository.CreateRole(ctx, payload)
}

func (u *RolesUsecase) CreatePermission(ctx context.Context, payload *roles.CreatePermissionPayload) (*roles.Permissions, error) {
	if payload.Code == "" || payload.Module == "" {
		return nil, fmt.Errorf("permission code and module are required")
	}
	return u.repository.CreatePermission(ctx, payload)
}

func (u *RolesUsecase) AssignPermissionToRole(ctx context.Context, payload *roles.CreateRolePermissionPayload) (*roles.RolePermissions, error) {
	if err := u.repository.ValidateRoleExists(ctx, payload.RoleID); err != nil {
		return nil, err
	}
	return u.repository.CreateRolePermission(ctx, payload)
}

func (u *RolesUsecase) GetAllRoles(ctx context.Context, limit, offset int) ([]roles.RoleWithPermissions, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return u.repository.GetAllRolesWithPermissions(ctx, limit, offset)
}

func (u *RolesUsecase) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*roles.RoleWithPermissions, error) {
	return u.repository.GetRoleWithPermissions(ctx, roleID)
}
