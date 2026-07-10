package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AssignRole adds a role to a staff member and writes an 'assigned' history entry.
// Both writes happen in a single transaction.
func (s *StaffRepository) AssignRole(ctx context.Context, staffID, roleID, performedBy uuid.UUID, reason *string) (*staff.StaffRoleResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert into staff_roles and resolve role name in one round-trip.
	// ON CONFLICT DO NOTHING returns no row if the role was already assigned,
	// which causes pgx.ErrNoRows — we surface that as a clear error.
	var roleName string
	err = tx.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO staff_roles (staff_id, role_id, assigned_by)
			VALUES ($1, $2, $3)
			ON CONFLICT (staff_id, role_id) DO NOTHING
			RETURNING role_id
		)
		SELECT r.name FROM ins JOIN roles r ON r.id = ins.role_id
	`, staffID, roleID, performedBy).Scan(&roleName)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role already assigned to this staff member")
		}
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	if err := insertRoleHistory(ctx, tx, staffID, roleID, performedBy, "assigned", reason); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &staff.StaffRoleResponse{
		StaffID:  staffID,
		RoleID:   roleID,
		RoleName: roleName,
		Action:   "assigned",
	}, nil
}

// RevokeRole removes a role from a staff member and writes a 'revoked' history entry.
func (s *StaffRepository) RevokeRole(ctx context.Context, staffID, roleID, performedBy uuid.UUID, reason *string) (*staff.StaffRoleResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Fetch role name and confirm assignment exists before deleting
	var roleName string
	if err := tx.QueryRow(ctx, `
		SELECT r.name
		FROM roles r
		JOIN staff_roles sr ON r.id = sr.role_id
		WHERE sr.staff_id = $1 AND sr.role_id = $2
	`, staffID, roleID).Scan(&roleName); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not assigned to this staff member")
		}
		return nil, fmt.Errorf("failed to verify role assignment: %w", err)
	}

	if _, err := tx.Exec(ctx,
		`DELETE FROM staff_roles WHERE staff_id = $1 AND role_id = $2`,
		staffID, roleID,
	); err != nil {
		return nil, fmt.Errorf("failed to revoke role: %w", err)
	}

	if err := insertRoleHistory(ctx, tx, staffID, roleID, performedBy, "revoked", reason); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &staff.StaffRoleResponse{
		StaffID:  staffID,
		RoleID:   roleID,
		RoleName: roleName,
		Action:   "revoked",
	}, nil
}

// UpdateRole atomically revokes the old role and assigns the new one.
// Writes two history entries in the same transaction.
func (s *StaffRepository) UpdateRole(ctx context.Context, staffID, performedBy uuid.UUID, payload *staff.UpdateRolePayload) (*staff.StaffRoleResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Verify old role is actually assigned
	var oldRoleName string
	if err := tx.QueryRow(ctx, `
		SELECT r.name
		FROM roles r
		JOIN staff_roles sr ON r.id = sr.role_id
		WHERE sr.staff_id = $1 AND sr.role_id = $2
	`, staffID, payload.OldRoleID).Scan(&oldRoleName); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("old role not assigned to this staff member")
		}
		return nil, fmt.Errorf("failed to verify old role: %w", err)
	}

	// Revoke old role
	if _, err := tx.Exec(ctx,
		`DELETE FROM staff_roles WHERE staff_id = $1 AND role_id = $2`,
		staffID, payload.OldRoleID,
	); err != nil {
		return nil, fmt.Errorf("failed to remove old role: %w", err)
	}
	if err := insertRoleHistory(ctx, tx, staffID, payload.OldRoleID, performedBy, "revoked", payload.Reason); err != nil {
		return nil, err
	}

	// Assign new role and resolve its name
	var newRoleName string
	if err := tx.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO staff_roles (staff_id, role_id, assigned_by)
			VALUES ($1, $2, $3)
			RETURNING role_id
		)
		SELECT r.name FROM ins JOIN roles r ON r.id = ins.role_id
	`, staffID, payload.NewRoleID, performedBy).Scan(&newRoleName); err != nil {
		return nil, fmt.Errorf("failed to assign new role: %w", err)
	}
	if err := insertRoleHistory(ctx, tx, staffID, payload.NewRoleID, performedBy, "assigned", payload.Reason); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &staff.StaffRoleResponse{
		StaffID:  staffID,
		RoleID:   payload.NewRoleID,
		RoleName: newRoleName,
		Action:   "assigned",
	}, nil
}

// GetRoleHistory returns the full role change log for a staff member, newest first.
func (s *StaffRepository) GetRoleHistory(ctx context.Context, staffID uuid.UUID) ([]staff.StaffRoleHistoryEntry, error) {
	stmt := `
		SELECT h.id, h.role_id, r.name, h.action, h.performed_by, h.reason, h.created_at
		FROM staff_role_history h
		JOIN roles r ON r.id = h.role_id
		WHERE h.staff_id = $1
		ORDER BY h.created_at DESC
	`

	rows, err := s.db.Query(ctx, stmt, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch role history: %w", err)
	}
	defer rows.Close()

	var list []staff.StaffRoleHistoryEntry
	for rows.Next() {
		var e staff.StaffRoleHistoryEntry
		if err := rows.Scan(&e.ID, &e.RoleID, &e.RoleName, &e.Action, &e.PerformedBy, &e.Reason, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}
		list = append(list, e)
	}
	return list, nil
}

// insertRoleHistory writes one row into staff_role_history within an open transaction.
func insertRoleHistory(ctx context.Context, tx pgx.Tx, staffID, roleID, performedBy uuid.UUID, action string, reason *string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO staff_role_history (staff_id, role_id, action, performed_by, reason)
		VALUES ($1, $2, $3, $4, $5)
	`, staffID, roleID, action, performedBy, reason)
	if err != nil {
		return fmt.Errorf("failed to write role history: %w", err)
	}
	return nil
}
