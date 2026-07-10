package repository

import (
	"context"
	"fmt"

	"github.com/bloansbook/bloansbook-api/internal/models/roles"
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
	return &StaffRepository{db: db, config: config}
}

// staffSelectQuery is the shared SELECT used by GetStaffByID and GetStaffByStaffID.
const staffSelectQuery = `
	SELECT
		s.id, s.staff_id, s.password_hash, s.first_name, s.last_name, s.email, s.phone, s.address,
		s.date_of_birth, s.date_of_hire, s.emergency_contact_name, s.emergency_contact_phone,
		s.bank_name, s.bank_account_number, s.bank_account_name, s.department, s.job_title,
		s.pay_type, s.base_salary, s.status, s.fired_at, s.has_login, s.supabase_uid,
		s.created_at, s.updated_at,
		creator.id AS creator_id, creator.staff_id AS creator_staff_id,
		creator.first_name AS creator_first_name, creator.last_name AS creator_last_name,
		creator.email AS creator_email, creator.phone AS creator_phone,
		creator.department AS creator_department, creator.job_title AS creator_job_title,
		creator.status AS creator_status
	FROM staff s
	LEFT JOIN staff creator ON s.created_by = creator.id
`

func (s *StaffRepository) CountStaff(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM staff`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}
	return count, nil
}

// CreateStaff inserts a new staff member. Pass a non-nil fixedID to use a specific UUID,
// or nil to let the database generate one.
func (s *StaffRepository) CreateStaff(ctx context.Context, fixedID *uuid.UUID, createdBy uuid.UUID, payload *staff.CreateStaffPayload, credentials staff.Credentials) (*staff.StaffCreate, error) {
	var exists bool
	if err := s.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM staff WHERE staff_id = $1)`,
		credentials.StaffID,
	).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check staff_id: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("staff_id %s already exists", credentials.StaffID)
	}

	var (
		stmt string
		args pgx.NamedArgs
	)

	baseArgs := pgx.NamedArgs{
		"staff_id":                credentials.StaffID,
		"password_hash":           credentials.Password,
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
		"supabase_uid":            payload.SupabaseUID,
	}

	if fixedID != nil {
		baseArgs["id"] = *fixedID
		stmt = `
			INSERT INTO staff (
				id, staff_id, password_hash, first_name, last_name, email, phone, address,
				date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone,
				bank_name, bank_account_number, bank_account_name, department, job_title,
				pay_type, base_salary, status, created_by, has_login
			) VALUES (
				@id, @staff_id, @password_hash, @first_name, @last_name, @email, @phone, @address,
				@date_of_birth, @date_of_hire, @emergency_contact_name, @emergency_contact_phone,
				@bank_name, @bank_account_number, @bank_account_name, @department, @job_title,
				@pay_type, @base_salary, @status, @created_by, @has_login
			) RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, created_at, updated_at, status
		`
	} else {
		stmt = `
			INSERT INTO staff (
				staff_id, password_hash, first_name, last_name, email, phone, address,
				date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone,
				bank_name, bank_account_number, bank_account_name, department, job_title,
				pay_type, base_salary, status, created_by, supabase_uid, has_login
			) VALUES (
				@staff_id, @password_hash, @first_name, @last_name, @email, @phone, @address,
				@date_of_birth, @date_of_hire, @emergency_contact_name, @emergency_contact_phone,
				@bank_name, @bank_account_number, @bank_account_name, @department, @job_title,
				@pay_type, @base_salary, @status, @created_by, @supabase_uid, @has_login
			) RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, created_at, updated_at, status
		`
	}
	args = baseArgs

	var m staff.StaffCreate
	err := s.db.QueryRow(ctx, stmt, args).Scan(
		&m.ID,
		&m.StaffID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Phone,
		&m.Department,
		&m.JobTitle,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}
	return &m, nil
}

func (s *StaffRepository) GetStaffByID(ctx context.Context, id uuid.UUID) (*staff.Staff, error) {
	stmt := staffSelectQuery + `WHERE s.id = $1`

	var m staff.Staff
	err := s.db.QueryRow(ctx, stmt, id).Scan(
		&m.ID, &m.StaffID, &m.Password,
		&m.FirstName, &m.LastName, &m.Email, &m.Phone, &m.Address,
		&m.DateOfBirth, &m.DateOfHire, &m.EmergencyContactName, &m.EmergencyContactPhone,
		&m.BankName, &m.BankAccountNumber, &m.BankAccountName,
		&m.Department, &m.JobTitle, &m.PayType, &m.BaseSalary,
		&m.Status, &m.FiredAt, &m.HasLogin, &m.SupabaseUID,
		&m.CreatedAt, &m.UpdatedAt,
		&m.CreatorID, &m.CreatorStaffID,
		&m.CreatorFirstName, &m.CreatorLastName,
		&m.CreatorEmail, &m.CreatorPhone,
		&m.CreatorDepartment, &m.CreatorJobTitle, &m.CreatorStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return &m, nil
}

func (s *StaffRepository) GetStaffByStaffID(ctx context.Context, staffID string) (*staff.Staff, error) {
	stmt := staffSelectQuery + `WHERE s.staff_id = $1`

	var m staff.Staff
	err := s.db.QueryRow(ctx, stmt, staffID).Scan(
		&m.ID, &m.StaffID, &m.Password,
		&m.FirstName, &m.LastName, &m.Email, &m.Phone, &m.Address,
		&m.DateOfBirth, &m.DateOfHire, &m.EmergencyContactName, &m.EmergencyContactPhone,
		&m.BankName, &m.BankAccountNumber, &m.BankAccountName,
		&m.Department, &m.JobTitle, &m.PayType, &m.BaseSalary,
		&m.Status, &m.FiredAt, &m.HasLogin, &m.SupabaseUID,
		&m.CreatedAt, &m.UpdatedAt,
		&m.CreatorID, &m.CreatorStaffID,
		&m.CreatorFirstName, &m.CreatorLastName,
		&m.CreatorEmail, &m.CreatorPhone,
		&m.CreatorDepartment, &m.CreatorJobTitle, &m.CreatorStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return &m, nil
}

func (s *StaffRepository) GetAllStaff(ctx context.Context, limit, offset int) ([]staff.Staff, error) {
	stmt := staffSelectQuery + `ORDER BY s.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all staff: %w", err)
	}
	defer rows.Close()

	var list []staff.Staff
	for rows.Next() {
		var m staff.Staff
		if err := rows.Scan(
			&m.ID, &m.StaffID, &m.Password,
			&m.FirstName, &m.LastName, &m.Email, &m.Phone, &m.Address,
			&m.DateOfBirth, &m.DateOfHire, &m.EmergencyContactName, &m.EmergencyContactPhone,
			&m.BankName, &m.BankAccountNumber, &m.BankAccountName,
			&m.Department, &m.JobTitle, &m.PayType, &m.BaseSalary,
			&m.Status, &m.FiredAt, &m.HasLogin, &m.SupabaseUID,
			&m.CreatedAt, &m.UpdatedAt,
			&m.CreatorID, &m.CreatorStaffID,
			&m.CreatorFirstName, &m.CreatorLastName,
			&m.CreatorEmail, &m.CreatorPhone,
			&m.CreatorDepartment, &m.CreatorJobTitle, &m.CreatorStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan staff row: %w", err)
		}
		list = append(list, m)
	}
	return list, nil
}

func (s *StaffRepository) UpdateStaff(ctx context.Context, id uuid.UUID, payload *staff.UpdateStaffPayload) (*staff.StaffCreate, error) {
	stmt := `
		UPDATE staff SET
			first_name              = COALESCE(@first_name, first_name),
			last_name               = COALESCE(@last_name, last_name),
			email                   = COALESCE(@email, email),
			phone                   = COALESCE(@phone, phone),
			address                 = COALESCE(@address, address),
			date_of_birth           = COALESCE(@date_of_birth, date_of_birth),
			emergency_contact_name  = COALESCE(@emergency_contact_name, emergency_contact_name),
			emergency_contact_phone = COALESCE(@emergency_contact_phone, emergency_contact_phone),
			bank_name               = COALESCE(@bank_name, bank_name),
			bank_account_number     = COALESCE(@bank_account_number, bank_account_number),
			bank_account_name       = COALESCE(@bank_account_name, bank_account_name),
			department              = COALESCE(@department, department),
			job_title               = COALESCE(@job_title, job_title),
			pay_type                = COALESCE(@pay_type, pay_type),
			base_salary             = COALESCE(@base_salary, base_salary),
			status                  = COALESCE(@status, status),
			updated_at              = NOW()
		WHERE id = @id
		RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, status, created_at, updated_at
	`

	var m staff.StaffCreate
	err := s.db.QueryRow(ctx, stmt, pgx.NamedArgs{
		"id":                      id,
		"first_name":              payload.FirstName,
		"last_name":               payload.LastName,
		"email":                   payload.Email,
		"phone":                   payload.Phone,
		"address":                 payload.Address,
		"date_of_birth":           payload.DateOfBirth,
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
	}).Scan(
		&m.ID, &m.StaffID,
		&m.FirstName, &m.LastName, &m.Email, &m.Phone,
		&m.Department, &m.JobTitle, &m.Status,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}
	return &m, nil
}

func (s *StaffRepository) GetStaffBySupabaseUID(ctx context.Context, supabaseUID string) (*staff.Staff, error) {
	stmt := staffSelectQuery + `WHERE s.supabase_uid = $1`

	var m staff.Staff
	err := s.db.QueryRow(ctx, stmt, supabaseUID).Scan(
		&m.ID, &m.StaffID, &m.Password,
		&m.FirstName, &m.LastName, &m.Email, &m.Phone, &m.Address,
		&m.DateOfBirth, &m.DateOfHire, &m.EmergencyContactName, &m.EmergencyContactPhone,
		&m.BankName, &m.BankAccountNumber, &m.BankAccountName,
		&m.Department, &m.JobTitle, &m.PayType, &m.BaseSalary,
		&m.Status, &m.FiredAt, &m.HasLogin, &m.SupabaseUID,
		&m.CreatedAt, &m.UpdatedAt,
		&m.CreatorID, &m.CreatorStaffID,
		&m.CreatorFirstName, &m.CreatorLastName,
		&m.CreatorEmail, &m.CreatorPhone,
		&m.CreatorDepartment, &m.CreatorJobTitle, &m.CreatorStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("staff not found: %w", err)
	}
	return &m, nil
}

// GetStaffPermissions returns the flat list of permission codes assigned to a
// staff member through their role(s). Used by the auth middleware.
func (s *StaffRepository) GetStaffPermissions(ctx context.Context, id uuid.UUID) ([]string, error) {
	stmt := `
		SELECT DISTINCT p.code
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN staff_roles sr ON rp.role_id = sr.role_id
		WHERE sr.staff_id = $1
		ORDER BY p.code
	`

	rows, err := s.db.Query(ctx, stmt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staff permissions: %w", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		codes = append(codes, code)
	}
	return codes, nil
}

func (s *StaffRepository) GetStaffRoles(ctx context.Context, id uuid.UUID) ([]roles.RoleSummary, error) {
	stmt := `
		SELECT r.id, r.name
		FROM roles r
		JOIN staff_roles sr ON r.id = sr.role_id
		WHERE sr.staff_id = $1
		ORDER BY r.name
	`

	rows, err := s.db.Query(ctx, stmt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staff roles: %w", err)
	}
	defer rows.Close()

	var list []roles.RoleSummary
	for rows.Next() {
		var r roles.RoleSummary
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		list = append(list, r)
	}
	return list, nil
}

// TODO: FIRE STAFF REPO
