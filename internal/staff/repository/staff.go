package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository struct {
	db     *pgxpool.Pool
	config *config.Config
}

func NewStaffRepository(db *pgxpool.Pool, config *config.Config) *StaffRepository {
	return &StaffRepository{
		config: config,
		db:     db,
	}
}

func (s *StaffRepository) CreateStaff(ctx context.Context, createdBy uuid.UUID, payload *staff.CreateStaffPayload) (*staff.Staff, error) {
	stmt := `
	INSERT INTO staff (
		staff_id, password_hash, first_name, last_name, email, phone, address, date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone, bank_name, bank_account_number, bank_account_name, department, job_title, pay_type, base_salary, status, created_by, has_login
	) VALUES (
	 	@staff_id, @password_hash, @first_name, @last_name, @email, @phone, @address, @date_of_birth, @date_of_hire, @emergency_contact_name, @emergency_contact_phone,
		@bank_name, @bank_account_number, @bank_account_name, @department, @job_title, @pay_type, @base_salary, @status, @created_by, @has_login
	) RETURNING id, staff_id, first_name, last_name, email, phone, address, date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone, bank_name, bank_account_number, bank_account_name, department, job_title, pay_type, base_salary, status, fired_at, has_login, superbase_uid, created_by, created_at, updated_at
	`

	var exists bool
	checkStmt := `SELECT EXISTS(SELECT 1 FROM staff WHERE staff_id = @staff_id)`
	err := s.db.QueryRow(ctx, checkStmt, pgx.NamedArgs{"staff_id": payload.StaffID}).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if staff exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("staff_id %s already exists", payload.StaffID)
	}

	rows, err := s.db.Query(ctx, stmt, pgx.NamedArgs{
		"staff_id":                payload.StaffID,
		"password_hash":           payload.Password,
		"first_name":              payload.FirstName,
		"last_name":               payload.LastName,
		"email":                   payload.Email,
		"phone":                   payload.Phone,
		"address":                 payload.Address,
		"date_of_birth":           payload.DateOfBirth,
		"date_of_hire":            payload.DateOfHire,
		"emergency_contact_name":  payload.EmergencyContactName,
		"emergency_contact_phone": payload.EmergencyContactPhone,
		"bank_name":               payload.BankName,
		"bank_account_number":     payload.BankAccountNumber,
		"bank_account_name":       payload.BankAccountName,
		"department":              payload.Department,
		"job_title":               payload.JobTitle,
		"pay_type":                payload.PayType,
		"base_salary":             payload.BaseSalary,
		"status":                  payload.Status,
		"created_by":              createdBy,
		"has_login":               payload.HasLogin,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute query for creating staff: %w", err)
	}

	staff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.Staff])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created staff: %w", err)
	}

	return &staff, nil
}
