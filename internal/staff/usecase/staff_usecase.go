package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bloansbook/bloansbook-api/internal/auth"
	a "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/email"
	"github.com/bloansbook/bloansbook-api/pkg/idgen"
	"github.com/bloansbook/bloansbook-api/pkg/password"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type StaffUsecase struct {
	authRepo     *a.AuthRepository
	emailService *email.Service
	repository   *repository.StaffRepository
	db           *pgxpool.Pool
	config       *config.Config
}

func NewStaffUsecase(db *pgxpool.Pool, auth *a.AuthRepository, repo *repository.StaffRepository, config *config.Config) *StaffUsecase {
	return &StaffUsecase{
		repository:   repo,
		db:           db,
		config:       config,
		authRepo:     auth,
		emailService: email.NewService(config.Resend.From),
	}
}

func (u *StaffUsecase) GetStaffCount(ctx context.Context) (int, error) {
	return u.repository.CountStaff(ctx)
}

func (u *StaffUsecase) CreateStaff(ctx context.Context, createdBy uuid.UUID, payload *staff.CreateStaffPayload) (*staff.CreateStaffResponse, error) {
	count, err := u.repository.CountStaff(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count staff: %w", err)
	}

	var creator uuid.UUID
	if count == 0 {
		sys, err := u.getOrCreateSystemStaff(ctx)
		if err != nil {
			return nil, err
		}
		creator = sys.ID
	} else {
		if createdBy == uuid.Nil {
			return nil, fmt.Errorf("authentication required: createdBy is missing")
		}
		creator = createdBy
	}

	staffID, err := idgen.GenerateSequentialID(ctx, u.db, "staff", "staff_id", u.config.IDGen.StaffPrefix)
	if err != nil {
		return nil, err
	}

	pass, err := password.GeneratePassword()
	if err != nil {
		return nil, err
	}

	if payload.HasLogin {
		res, err := u.authRepo.AdminCreateUser(ctx, auth.LoginDTO{StaffID: staffID, Password: pass})
		if err != nil {
			return nil, err
		}
		uid := res.User.ID.String()
		payload.SupabaseUID = &uid
	}

	hashedPassword, err := password.HashPassword(pass)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	credentials := staff.Credentials{StaffID: staffID, Password: hashedPassword}

	m, err := u.repository.CreateStaff(ctx, nil, creator, payload, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	if err := u.emailService.SendStaffWelcomeEmail(ctx, *m.Email, m.StaffID, pass, m.FirstName, m.LastName); err != nil {
		return nil, fmt.Errorf("failed to send welcome email: %w", err)
	}

	return &staff.CreateStaffResponse{
		ID:        m.ID,
		Staff:     m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (u *StaffUsecase) GetStaffByID(ctx context.Context, id uuid.UUID) (*staff.StaffDTO, error) {
	m, err := u.repository.GetStaffByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}

	staffRoles, err := u.repository.GetStaffRoles(ctx, m.ID)
	if err != nil {
		fmt.Printf("warning: failed to fetch roles for staff %s: %v\n", m.ID, err)
		staffRoles = nil
	}

	dto := m.ToDTO(staffRoles)
	return &dto, nil
}

func (u *StaffUsecase) GetAllStaff(ctx context.Context, limit, offset int) ([]staff.StaffDTO, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	members, err := u.repository.GetAllStaff(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	dtos := make([]staff.StaffDTO, len(members))
	for i, m := range members {
		dtos[i] = m.ToDTO(nil)
	}
	return dtos, nil
}

func (u *StaffUsecase) UpdateStaff(ctx context.Context, id uuid.UUID, payload *staff.UpdateStaffPayload) (*staff.UpdateStaffResponse, error) {
	if _, err := u.repository.GetStaffByID(ctx, id); err != nil {
		return nil, fmt.Errorf("cannot update non-existent staff: %w", err)
	}

	m, err := u.repository.UpdateStaff(ctx, id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return &staff.UpdateStaffResponse{
		ID:        m.ID,
		Staff:     m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (u *StaffUsecase) getOrCreateSystemStaff(ctx context.Context) (*staff.CreateStaffResponse, error) {
	sys, err := u.repository.GetStaffByStaffID(ctx, "SYSTEM-0000")
	if err == nil {
		return &staff.CreateStaffResponse{
			ID:        sys.ID,
			Staff:     sys.ToSummary(),
			CreatedAt: sys.CreatedAt,
			UpdatedAt: sys.UpdatedAt,
		}, nil
	}

	if !strings.Contains(err.Error(), "not found") {
		return nil, err
	}

	credentials := staff.Credentials{
		StaffID:  "SYSTEM-0000",
		Password: "SYSTEM_ACCOUNT_DO_NOT_USE_$$$",
	}

	systemPayload := &staff.CreateStaffPayload{
		FirstName:  "System",
		LastName:   "Account",
		Department: "system",
		JobTitle:   "system_account",
		PayType:    "monthly",
		BaseSalary: decimal.NewFromInt(0),
		Status:     staff.StaffStatusActive,
		HasLogin:   false,
		DateOfHire: time.Now(),
	}

	systemUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	m, err := u.repository.CreateStaff(ctx, &systemUUID, systemUUID, systemPayload, credentials)
	if err != nil {
		return nil, err
	}

	return &staff.CreateStaffResponse{
		ID:        m.ID,
		Staff:     m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// --- Staff Role Methods ---

func (u *StaffUsecase) AssignRole(ctx context.Context, staffID, roleID, performedBy uuid.UUID, reason *string) (*staff.StaffRoleResponse, error) {
	if _, err := u.repository.GetStaffByID(ctx, staffID); err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return u.repository.AssignRole(ctx, staffID, roleID, performedBy, reason)
}

func (u *StaffUsecase) RevokeRole(ctx context.Context, staffID, roleID, performedBy uuid.UUID, reason *string) (*staff.StaffRoleResponse, error) {
	if _, err := u.repository.GetStaffByID(ctx, staffID); err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return u.repository.RevokeRole(ctx, staffID, roleID, performedBy, reason)
}

func (u *StaffUsecase) UpdateRole(ctx context.Context, staffID, performedBy uuid.UUID, payload *staff.UpdateRolePayload) (*staff.StaffRoleResponse, error) {
	if _, err := u.repository.GetStaffByID(ctx, staffID); err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return u.repository.UpdateRole(ctx, staffID, performedBy, payload)
}

func (u *StaffUsecase) GetRoleHistory(ctx context.Context, staffID uuid.UUID) ([]staff.StaffRoleHistoryEntry, error) {
	if _, err := u.repository.GetStaffByID(ctx, staffID); err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return u.repository.GetRoleHistory(ctx, staffID)
}
