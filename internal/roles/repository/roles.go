package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RolesRepository struct {
	db     *pgxpool.Pool
	config *config.Config
}

func NewRolesRepository(db *pgxpool.Pool, config *config.Config) *RolesRepository {
	return &RolesRepository{
		db:     db,
		config: config,
	}
}

func (r *RolesRepository) validateRoleExists(ctx context.Context, roleID uuid.UUID) error {
	stmt := `SELECT EXISTS(SELECT 1 FROM roles WHERE id = @role_id)`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"role_id": roleID,
	})
	if err != nil {
		return fmt.Errorf("failed to check role existence: %w", err)
	}

	var exists bool
	_, err = pgx.ForEachRow(rows, func(row pgx.CollectableRow) error {
		return row.Scan(&exists)
	})
	if err != nil {
		return fmt.Errorf("failed to scan existence result: %w", err)
	}

	if !exists {
		return fmt.Errorf("role with ID %s does not exist", roleID.String())
	}

	return nil
}

func (r *RolesRepository) CreateRole(ctx context.Context, payload *roles.CreateRolePayload) (*roles.Roles, error) {
	stmt := `
		INSERT INTO roles (name, description, is_system)
		VALUES (@name, @description, @is_system)
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"name":        payload.Name,
		"description": payload.Description,
		"is_system":   payload.IsSystem,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for creating role: %w", err)
	}

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roles.Roles])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created role: %w", err)
	}

	return &role, nil
}

func (r *RolesRepository) CreatePermission(ctx context.Context, payload *roles.CreatePermissionPayload) (*roles.Permissions, error) {
	stmt := `
		INSERT INTO permissions (code, module, description)
		VALUES (@code, @module, @description)
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"code":        payload.Code,
		"module":      payload.Module,
		"description": payload.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for creating permission: %w", err)
	}

	permission, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roles.Permissions])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created permission: %w", err)
	}

	return &permission, nil
}

func (r *RolesRepository) CreateRolePermission(ctx context.Context, payload *roles.CreateRolePermissionPayload) (*roles.RolePermissions, error) {
	stmt := `
		WITH inserted AS (
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES (@role_id, @permission_id)
			RETURNING role_id, permission_id, created_at
		)
		SELECT
			inserted.created_at,
			roles.id AS role_id,
			permissions.id AS permission_id,
			roles.name AS role_name,
			permissions.code AS permission_code,
			permissions.module AS permission_module
		FROM inserted
		JOIN roles ON roles.id = inserted.role_id
		JOIN permissions ON permissions.id = inserted.permission_id
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"role_id":       payload.RoleID,
		"permission_id": payload.PermissionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for creating role permission: %w", err)
	}

	rolePermission, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roles.RolePermissions])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created role permission: %w", err)
	}

	return &rolePermission, nil
}

func (r *RolesRepository) GetAllRoles(ctx context.Context) ([]roles.Roles, error) {
	stmt := `SELECT * FROM roles`

	rows, err := r.db.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting all roles: %w", err)
	}

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByName[roles.Roles])
	if err != nil {
		return nil, fmt.Errorf("failed to collect roles: %w", err)
	}

	return roles, nil
}

func (r *RolesRepository) GetRoleByName(ctx context.Context, name string) (*roles.Roles, error) {
	stmt := `SELECT * FROM roles WHERE name = @name`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"name": name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting role by name: %w", err)
	}

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roles.Roles])
	if err != nil {
		return nil, fmt.Errorf("failed to collect role by name: %w", err)
	}

	return &role, nil
}

func (r *RolesRepository) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*roles.RoleWithPermissions, error) {
	stmt := `
	SELECT
		r.id,
		r.name,
		COALESCE(
			json_agg(
				json_build_object(
					'id', p.id,
					'code', p.code,
					'module', p.module
				) ORDER BY p.code
			) FILTER (WHERE p.id IS NOT NULL),
			'[]'::json
		) AS permissions
	FROM roles r
	LEFT JOIN role_permissions rp ON r.id = rp.role_id
	LEFT JOIN permissions p ON rp.permission_id = p.id
	WHERE r.id = @role_id
	GROUP BY r.id, r.name
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"role_id": roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting role with permissions: %w", err)
	}

	roleWithPermissions, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roles.RoleWithPermissions])
	if err != nil {
		return nil, fmt.Errorf("failed to collect role with permissions: %w", err)
	}

	return &roleWithPermissions, nil
}

func (r *RolesRepository) GetAllRolesWithPermissions(ctx context.Context, limit, offset int) ([]roles.RoleWithPermissions, error) {
	stmt := `
	SELECT r.id,
	r.name,
	COALESCE(
		json_agg(
			json_build_object(
				'id', p.id,
				'code', p.code,
				'module', p.module
			) ORDER BY p.code
		) FILTER (WHERE p.id IS NOT NULL),
		 '[]'::json
	) AS permissions
	FROM roles r
	LEFT JOIN role_permissions rp ON r.id = rp.role_id
	LEFT JOIN permissions p ON rp.permission_id = p.id
	GROUP BY r.id, r.name
	ORDER BY r.name
	LIMIT @limit OFFSET @offset
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting all roles with permissions: %w", err)
	}

	rolesWithPermissions, err := pgx.CollectRows(rows, pgx.RowToStructByName[roles.RoleWithPermissions])
	if err != nil {
		return nil, fmt.Errorf("failed to collect roles with permissions: %w", err)
	}

	return rolesWithPermissions, nil
}
