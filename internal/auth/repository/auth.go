package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	s "github.com/mrehanabbasi/supabase-auth-go"
	"github.com/mrehanabbasi/supabase-auth-go/types"
)

type AuthRepository struct {
	client      s.Client
	adminClient s.Client
	db          *pgxpool.Pool
}

func NewAuthRepository(cfg *config.Config, db *pgxpool.Pool) *AuthRepository {
	client := s.New(cfg.Supabase.ProjectRef, cfg.Supabase.AnonKey)
	return &AuthRepository{
		client:      client,
		adminClient: client.WithToken(cfg.Supabase.ServiceRoleKey),
		db:          db,
	}
}

func (a *AuthRepository) SignInWithPassword(ctx context.Context, payload auth.LoginDTO) (*types.TokenResponse, error) {
	resp, err := a.client.Token(ctx, types.TokenRequest{
		GrantType: "password",
		Email:     fmt.Sprintf("%s@bloansbook.local", payload.StaffID),
		Password:  payload.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}
	return resp, nil
}

func (a *AuthRepository) RefreshToken(ctx context.Context, refreshToken string) (*types.TokenResponse, error) {
	resp, err := a.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return resp, nil
}

func (a *AuthRepository) AdminCreateUser(ctx context.Context, payload auth.LoginDTO) (*types.AdminCreateUserResponse, error) {
	resp, err := a.adminClient.AdminCreateUser(ctx, types.AdminCreateUserRequest{
		Email:        fmt.Sprintf("%s@bloansbook.local", payload.StaffID),
		Password:     &payload.Password,
		EmailConfirm: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auth user: %w", err)
	}
	return resp, nil
}

func (a *AuthRepository) GetUser(ctx context.Context, accessToken string) (*types.UserResponse, error) {
	user, err := a.client.WithToken(accessToken).GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (a *AuthRepository) GetMe(ctx context.Context, userID uuid.UUID) (*auth.ProfileDTO, error) {
	stmt := `
		WITH module_perms AS (
			SELECT
				rp.role_id,
				p.module,
				json_agg(p.code ORDER BY p.code) AS codes
			FROM role_permissions rp
			JOIN permissions p ON p.id = rp.permission_id
			GROUP BY rp.role_id, p.module
			),
			role_perms AS (
				-- Step 2: group modules (with their codes) into one array per role
				SELECT
					role_id,
					json_agg(
						jsonb_build_object('module', module, 'codes', codes)
						ORDER BY module
					) AS permissions
				FROM module_perms
				GROUP BY role_id
			)
			SELECT
				s.id,
				s.staff_id,
				s.first_name,
				s.last_name,
				s.email,
				s.phone,
				s.address,
				s.date_of_birth,
				s.date_of_hire,
				s.emergency_contact_name,
				s.emergency_contact_phone,
				s.bank_name,
				s.bank_account_number,
				s.bank_account_name,
				s.department,
				s.job_title,
				s.pay_type,
				s.base_salary,
				s.status,
				s.fired_at,
				s.created_at,
				s.updated_at,
			COALESCE(
				json_agg(
					jsonb_build_object(
						'id', r.id,
						'name', r.name,
						'permissions', COALESCE(rpm.permissions, '[]'::json)
					)
				) FILTER (WHERE r.id IS NOT NULL),
				'[]'
			) AS roles
		FROM staff s
		LEFT JOIN staff_roles sr ON sr.staff_id = s.id
		LEFT JOIN roles r ON r.id = sr.role_id
		LEFT JOIN role_perms rpm ON rpm.role_id = r.id
		WHERE s.id = $1
		GROUP BY s.id
	`
	var me auth.ProfileDTO
	err := a.db.QueryRow(ctx, stmt, userID).Scan(
		&me.Staff.ID,
		&me.Staff.StaffID,
		&me.Staff.FirstName,
		&me.Staff.LastName,
		&me.Staff.Email,
		&me.Staff.Phone,
		&me.Staff.Address,
		&me.Staff.DateOfBirth,
		&me.Staff.DateOfHire,
		&me.Staff.EmergencyContactName,
		&me.Staff.EmergencyContactPhone,
		&me.Staff.BankName,
		&me.Staff.BankAccountNumber,
		&me.Staff.BankAccountName,
		&me.Staff.Department,
		&me.Staff.JobTitle,
		&me.Staff.PayType,
		&me.Staff.BaseSalary,
		&me.Staff.Status,
		&me.Staff.FiredAt,
		&me.Staff.CreatedAt,
		&me.Staff.UpdatedAt,
		&me.Roles,
	)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &me, nil
}
