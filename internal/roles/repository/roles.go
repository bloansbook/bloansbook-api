package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RolesRepository struct {
	db     *pgxpool.Pool
	config *config.Config
}

func NewRolesRepository(db *pgxpool.Pool, config *config.Config) *RolesRepository {
	return &RolesRepository{db: db, config: config}
}

func (r *RolesRepository) ValidateRoleExists(ctx context.Context, roleID uuid.UUID) error {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)`,
		roleID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check role: %w", err)
	}
	if !exists {
		return fmt.Errorf("role %s does not exist", roleID)
	}
	return nil
}

func (r *RolesRepository) CreateRole(ctx context.Context, payload *roles.CreateRolePayload) (*roles.Roles, error) {
	stmt := `
		INSERT INTO roles (name, description, is_system)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, is_system, created_at, updated_at
	`

	var m roles.Roles
	err := r.db.QueryRow(ctx, stmt, payload.Name, payload.Description, payload.IsSystem).Scan(
		&m.ID, &m.Name, &m.Description, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return &m, nil
}

func (r *RolesRepository) CreatePermission(ctx context.Context, payload *roles.CreatePermissionPayload) (*roles.Permissions, error) {
	stmt := `
		INSERT INTO permissions (code, module, description)
		VALUES ($1, $2, $3)
		RETURNING id, code, module, description, created_at
	`

	var m roles.Permissions
	err := r.db.QueryRow(ctx, stmt, payload.Code, payload.Module, payload.Description).Scan(
		&m.ID, &m.Code, &m.Module, &m.Description, &m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}
	return &m, nil
}

func (r *RolesRepository) CreateRolePermission(ctx context.Context, payload *roles.CreateRolePermissionPayload) (*roles.RolePermissions, error) {
	stmt := `
		WITH inserted AS (
			INSERT INTO role_permissions (role_id, permission_id)
			VALUES ($1, $2)
			RETURNING role_id, permission_id, created_at
		)
		SELECT
			inserted.created_at,
			roles.id        AS role_id,
			permissions.id  AS permission_id,
			roles.name      AS role_name,
			permissions.code   AS permission_code,
			permissions.module AS permission_module
		FROM inserted
		JOIN roles       ON roles.id       = inserted.role_id
		JOIN permissions ON permissions.id = inserted.permission_id
	`

	var m roles.RolePermissions
	err := r.db.QueryRow(ctx, stmt, payload.RoleID, payload.PermissionID).Scan(
		&m.CreatedAt,
		&m.RoleID, &m.PermissionID,
		&m.Role, &m.PermissionCode, &m.PermissionModule,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assign permission to role: %w", err)
	}
	return &m, nil
}

func (r *RolesRepository) GetRoleByName(ctx context.Context, name string) (*roles.Roles, error) {
	stmt := `SELECT id, name, description, is_system, created_at, updated_at FROM roles WHERE name = $1`

	var m roles.Roles
	err := r.db.QueryRow(ctx, stmt, name).Scan(
		&m.ID, &m.Name, &m.Description, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &m, nil
}

func (r *RolesRepository) GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (*roles.RoleWithPermissions, error) {
	stmt := `
		SELECT
			r.id,
			r.name,
			COALESCE(
				json_agg(
					json_build_object('id', p.id, 'code', p.code, 'module', p.module)
					ORDER BY p.code
				) FILTER (WHERE p.id IS NOT NULL),
				'[]'::json
			) AS permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1
		GROUP BY r.id, r.name
	`

	var m roles.RoleWithPermissions
	err := r.db.QueryRow(ctx, stmt, roleID).Scan(&m.ID, &m.Name, &m.Permissions)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return &m, nil
}

func (r *RolesRepository) GetAllRolesWithPermissions(ctx context.Context, limit, offset int) ([]roles.RoleWithPermissions, error) {
	stmt := `
		SELECT
			r.id,
			r.name,
			COALESCE(
				json_agg(
					json_build_object('id', p.id, 'code', p.code, 'module', p.module)
					ORDER BY p.code
				) FILTER (WHERE p.id IS NOT NULL),
				'[]'::json
			) AS permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		GROUP BY r.id, r.name
		ORDER BY r.name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	var list []roles.RoleWithPermissions
	for rows.Next() {
		var m roles.RoleWithPermissions
		if err := rows.Scan(&m.ID, &m.Name, &m.Permissions); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *RolesRepository) GetAllRoles(ctx context.Context) ([]roles.Roles, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, description, is_system, created_at, updated_at FROM roles`)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	var list []roles.Roles
	for rows.Next() {
		var m roles.Roles
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		list = append(list, m)
	}
	return list, nil
}


