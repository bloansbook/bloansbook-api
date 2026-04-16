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
	return &RolesUsecase{
		repository: repo,
	}
}

// CreateRole creates a new role with the given payload
func (u *RolesUsecase) CreateRole(ctx context.Context, payload *roles.CreateRolePayload) (*roles.CreateRoleResponse, error) {
	if payload.Name == "" {
		return nil, fmt.Errorf("role name is required")
	}

	role, err := u.repository.CreateRole(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &roles.CreateRoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: &role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

// CreatePermission creates a new permission
func (u *RolesUsecase) CreatePermission(ctx context.Context, payload *roles.CreatePermissionPayload) (*roles.CreatePermissionResponse, error) {
	if payload.Code == "" || payload.Module == "" {
		return nil, fmt.Errorf("permission code and module are required")
	}

	permission, err := u.repository.CreatePermission(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &roles.CreatePermissionResponse{
		ID:          permission.ID,
		Code:        permission.Code,
		Module:      permission.Module,
		Description: &permission.Description,
		CreatedAt:   permission.CreatedAt,
	}, nil
}

// AssignPermissionToRole links a permission to a role
func (u *RolesUsecase) AssignPermissionToRole(ctx context.Context, payload *roles.CreateRolePermissionPayload) (*roles.CreateRolePermissionResponse, error) {
	exists, err := u.validateRoleExists(ctx, payload.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate role: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("role with ID %s does not exist", payload.RoleID.String())
	}

	rolePermission, err := u.repository.CreateRolePermission(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to assign permission to role: %w", err)
	}

	return &roles.CreateRolePermissionResponse{
		Role: roles.RoleSummary{
			ID:   rolePermission.RoleID,
			Name: rolePermission.Role,
		},
		Permission: roles.PermissionSummary{
			ID:     rolePermission.PermissionID,
			Code:   rolePermission.PermissionCode,
			Module: rolePermission.PermissionModule,
		},
		CreatedAt: rolePermission.CreatedAt,
	}, nil
}

// GetAllRoles returns all roles with their permissions
func (u *RolesUsecase) GetAllRoles(ctx context.Context, limit, offset int) ([]roles.RoleWithPermissions, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rolesWithPermissions, err := u.repository.GetAllRolesWithPermissions(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return rolesWithPermissions, nil
}

// GetRoleWithPermissions returns a single role with its permissions
func (u *RolesUsecase) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*roles.RoleWithPermissions, error) {
	roleWithPermissions, err := u.repository.GetRoleWithPermissions(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role with permissions: %w", err)
	}

	return roleWithPermissions, nil
}

// validateRoleExists checks if a role exists (moved from repository to usecase)
func (u *RolesUsecase) validateRoleExists(ctx context.Context, roleID uuid.UUID) (bool, error) {
	_, err := u.repository.GetRoleWithPermissions(ctx, roleID)
	if err != nil {
		return false, nil // Role doesn't exist
	}
	return true, nil
}
