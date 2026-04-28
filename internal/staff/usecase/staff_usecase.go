package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/idgen"
	"github.com/bloansbook/bloansbook-api/pkg/password"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// StaffUsecase handles staff business logic
type StaffUsecase struct {
	repository *repository.StaffRepository
	db         *pgxpool.Pool
	config     *config.Config
}

// NewStaffUsecase creates a new staff usecase
func NewStaffUsecase(db *pgxpool.Pool, repo *repository.StaffRepository, config *config.Config) *StaffUsecase {
	return &StaffUsecase{
		repository: repo,
		db:         db,
		config:     config,
	}
}

// Get Count of staff members
func (u *StaffUsecase) GetStaffCount(ctx context.Context) (int, error) {
	count, err := u.repository.CountStaff(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}

	return count, nil
}

// CreateStaff creates a new staff member
func (u *StaffUsecase) CreateStaff(ctx context.Context, createdBy uuid.UUID, payload *staff.CreateStaffPayload) (*staff.CreateStaffResponse, error) {
	// Count staff table
	count, err := u.repository.CountStaff(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count staff: %w", err)
	}

	var creator uuid.UUID
	if count == 0 {
		sysStaff, err := u.getOrCreateSystemStaff(ctx)
		if err != nil {
			return nil, err
		}
		creator = sysStaff.ID
	} else {
		if createdBy == uuid.Nil {
			return nil, fmt.Errorf("authentication required: user ID is required. %w", err)
		}
		creator = createdBy
	}
	// Generate staff id
	staffID, err := idgen.GenerateSequentialID(ctx, u.db, "staff", "staff_id", u.config.IDGen.StaffPrefix)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Generate password
	pass, err := password.GeneratePassword()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Hash password
	hashedPassword, err := password.HashPassword(pass)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	credentials := staff.Credentials{
		StaffID:  staffID,
		Password: hashedPassword,
	}

	// Validate required fields
	if payload.FirstName == "" || payload.LastName == "" {
		return nil, fmt.Errorf("first_name and last_name are required")
	}
	if payload.Email == nil || *payload.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if payload.Department == "" {
		return nil, fmt.Errorf("department is required")
	}
	if payload.JobTitle == "" {
		return nil, fmt.Errorf("job_title is required")
	}
	if payload.PayType == "" {
		return nil, fmt.Errorf("pay_type is required")
	}

	// Create staff
	staffMember, err := u.repository.CreateStaff(ctx, creator, payload, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return &staff.CreateStaffResponse{
		ID: staffMember.ID,
		Staff: staff.StaffSummary{
			StaffID:    staffMember.StaffID,
			FirstName:  staffMember.FirstName,
			LastName:   staffMember.LastName,
			Email:      staffMember.Email,
			Phone:      staffMember.Phone,
			Department: staffMember.Department,
			JobTitle:   staffMember.JobTitle,
			Status:     staffMember.Status,
		},
		CreatedAt: staffMember.CreatedAt,
		UpdatedAt: staffMember.UpdatedAt,
	}, nil

	// TODO: Implement Resend for email - to send staffID and password to staff
}

// GetStaffByID retrieves a staff member by ID
func (u *StaffUsecase) GetStaffByID(ctx context.Context, id uuid.UUID) (*staff.StaffDTO, error) {
	staffMember, err := u.repository.GetStaffByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return &staff.StaffDTO{
		ID:                    staffMember.ID,
		StaffID:               staffMember.StaffID,
		FirstName:             staffMember.FirstName,
		LastName:              staffMember.LastName,
		Email:                 staffMember.Email,
		Phone:                 staffMember.Phone,
		Address:               staffMember.Address,
		DateOfBirth:           staffMember.DateOfBirth,
		DateOfHire:            staffMember.DateOfHire,
		EmergencyContactName:  staffMember.EmergencyContactName,
		EmergencyContactPhone: staffMember.EmergencyContactPhone,
		BankName:              staffMember.BankName,
		BankAccountNumber:     staffMember.BankAccountNumber,
		BankAccountName:       staffMember.BankAccountName,
		Department:            staffMember.JobTitle,
		JobTitle:              staffMember.JobTitle,
		PayType:               staffMember.PayType,
		BaseSalary:            staffMember.BaseSalary,
		Status:                staffMember.Status,
		FiredAt:               staffMember.FiredAt,
		HasLogin:              staffMember.HasLogin,
		CreatedBy: staff.StaffSummary{
			StaffID:    staffMember.CreatorStaffID,
			FirstName:  staffMember.CreatorFirstName,
			LastName:   staffMember.CreatorLastName,
			Email:      staffMember.CreatorEmail,
			Phone:      staffMember.CreatorPhone,
			Department: staffMember.CreatorDepartment,
			JobTitle:   staffMember.CreatorJobTitle,
			Status:     staffMember.CreatorStatus,
		},
		CreatedAt: staffMember.CreatedAt,
		UpdatedAt: staffMember.UpdatedAt,
	}, nil
}

// GetAllStaff retrieves all staff members with pagination
func (u *StaffUsecase) GetAllStaff(ctx context.Context, limit, offset int) ([]staff.StaffDTO, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	staffMembers, err := u.repository.GetAllStaff(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff members: %w", err)
	}

	staffDTOs := make([]staff.StaffDTO, len(staffMembers))
	for i, m := range staffMembers {
		staffDTOs[i] = staff.StaffDTO{
			ID:                    m.ID,
			StaffID:               m.StaffID,
			FirstName:             m.FirstName,
			LastName:              m.LastName,
			Email:                 m.Email,
			Phone:                 m.Phone,
			Address:               m.Address,
			DateOfBirth:           m.DateOfBirth,
			DateOfHire:            m.DateOfHire,
			EmergencyContactName:  m.EmergencyContactName,
			EmergencyContactPhone: m.EmergencyContactPhone,
			BankName:              m.BankName,
			BankAccountNumber:     m.BankAccountNumber,
			BankAccountName:       m.BankAccountName,
			Department:            m.Department,
			JobTitle:              m.JobTitle,
			PayType:               m.PayType,
			BaseSalary:            m.BaseSalary,
			Status:                m.Status,
			FiredAt:               m.FiredAt,
			HasLogin:              m.HasLogin,
			CreatedBy: staff.StaffSummary{
				StaffID:    m.CreatorStaffID,
				FirstName:  m.CreatorFirstName,
				LastName:   m.CreatorLastName,
				Email:      m.CreatorEmail,
				Phone:      m.CreatorPhone,
				Department: m.CreatorDepartment,
				JobTitle:   m.CreatorJobTitle,
				Status:     m.CreatorStatus,
			},
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return staffDTOs, nil
}

// UpdateStaff updates staff member information
func (u *StaffUsecase) UpdateStaff(ctx context.Context, id uuid.UUID, payload *staff.UpdateStaffPayload) (*staff.UpdateStaffResponse, error) {
	// Validate staff exists
	_, err := u.repository.GetStaffByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot update non-existent staff: %w", err)
	}

	updatedStaff, err := u.repository.UpdateStaff(ctx, id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return &staff.UpdateStaffResponse{
		ID: updatedStaff.ID,
		Staff: staff.StaffSummary{
			StaffID:    updatedStaff.StaffID,
			FirstName:  updatedStaff.FirstName,
			LastName:   updatedStaff.LastName,
			Email:      updatedStaff.Email,
			Phone:      updatedStaff.Phone,
			Department: updatedStaff.Department,
			JobTitle:   updatedStaff.JobTitle,
			Status:     updatedStaff.Status,
		},
		CreatedAt: updatedStaff.CreatedAt,
		UpdatedAt: updatedStaff.UpdatedAt,
	}, nil
}

func (u *StaffUsecase) getOrCreateSystemStaff(ctx context.Context) (*staff.CreateStaffResponse, error) {
	sys, err := u.repository.GetStaffByStaffID(ctx, "SYSTEM-0000")
	if err == nil {
		return &staff.CreateStaffResponse{
			ID: sys.ID,
			Staff: staff.StaffSummary{
				StaffID:    sys.StaffID,
				Phone:      sys.Phone,
				Department: sys.Department,
				JobTitle:   sys.JobTitle,
				Status:     sys.Status,
			},
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
		FirstName:             "System",
		LastName:              "Account",
		Department:            "system",
		JobTitle:              "system_account",
		PayType:               "monthly",
		BaseSalary:            decimal.NewFromInt(0),
		Status:                staff.StaffStatusActive,
		HasLogin:              false,
		DateOfHire:            time.Now(),
		Email:                 nil,
		Phone:                 nil,
		Address:               nil,
		DateOfBirth:           nil,
		EmergencyContactName:  nil,
		EmergencyContactPhone: nil,
		BankName:              nil,
		BankAccountNumber:     nil,
		BankAccountName:       nil,
	}

	systemUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	createdStaff, err := u.repository.CreateStaffWithFixedID(ctx, systemUUID, systemUUID, systemPayload, credentials)
	if err != nil {
		return nil, err
	}

	return &staff.CreateStaffResponse{
		ID: createdStaff.ID,
		Staff: staff.StaffSummary{
			StaffID:    createdStaff.StaffID,
			Phone:      createdStaff.Phone,
			Department: createdStaff.Department,
			JobTitle:   createdStaff.JobTitle,
			Status:     createdStaff.Status,
		},
		CreatedAt: createdStaff.CreatedAt,
		UpdatedAt: createdStaff.UpdatedAt,
	}, nil
}
