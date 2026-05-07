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
	return &StaffRepository{
		config: config,
		db:     db,
	}
}

func (s *StaffRepository) CountStaff(ctx context.Context) (int, error) {
	stmt := `SELECT COUNT(*) FROM staff`

	var count int
	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to count staff: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to scan count: %w", err)
		}
	}

	return count, nil
}

func (s *StaffRepository) CreateStaff(ctx context.Context, createdBy uuid.UUID, payload *staff.CreateStaffPayload, credentials staff.Credentials) (*staff.StaffCreate, error) {
	stmt := `
		INSERT INTO staff (
			staff_id, password_hash, first_name, last_name, email, phone, address, date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone, bank_name, bank_account_number, bank_account_name, department, job_title, pay_type, base_salary, status, created_by, supabase_uid, has_login
		) VALUES (
			@staff_id, @password_hash, @first_name, @last_name, @email, @phone, @address, @date_of_birth, @date_of_hire, @emergency_contact_name, @emergency_contact_phone, @bank_name, @bank_account_number, @bank_account_name, @department, @job_title, @pay_type, @base_salary, @status, @created_by, @supabase_uid, @has_login
		) RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, created_at, updated_at, status
	`

	var exists bool
	checkStmt := `SELECT EXISTS(SELECT 1 FROM staff WHERE staff_id = $1)`
	err := s.db.QueryRow(ctx, checkStmt, credentials.StaffID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if staff exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("staff_id %s already exists", credentials.StaffID)
	}

	rows, err := s.db.Query(ctx, stmt, pgx.NamedArgs{
		"staff_id":                credentials.StaffID,
		"password_hash":           credentials.Password,
		"first_name":              payload.FirstName,
		"last_name":               payload.LastName,
		"email":                   &payload.Email,
		"phone":                   &payload.Phone,
		"address":                 &payload.Address,
		"date_of_birth":           &payload.DateOfBirth,
		"date_of_hire":            payload.DateOfHire,
		"emergency_contact_name":  &payload.EmergencyContactName,
		"emergency_contact_phone": &payload.EmergencyContactPhone,
		"bank_name":               &payload.BankName,
		"bank_account_number":     &payload.BankAccountNumber,
		"bank_account_name":       &payload.BankAccountName,
		"department":              payload.Department,
		"job_title":               payload.JobTitle,
		"pay_type":                payload.PayType,
		"base_salary":             payload.BaseSalary,
		"status":                  payload.Status,
		"created_by":              createdBy,
		"supabase_uid":            &payload.SupabaseUID,
		"has_login":               payload.HasLogin,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute query for creating staff: %w", err)
	}

	staff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.StaffCreate])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created staff: %w", err)
	}

	return &staff, nil
}

func (s *StaffRepository) CreateStaffWithFixedID(ctx context.Context, id, createdBy uuid.UUID, payload *staff.CreateStaffPayload, credentials staff.Credentials) (*staff.StaffCreate, error) {
	stmt := `
		INSERT INTO staff (
			id, staff_id, password_hash, first_name, last_name, email, phone, address, date_of_birth, date_of_hire, emergency_contact_name, emergency_contact_phone, bank_name, bank_account_number, bank_account_name, department, job_title, pay_type, base_salary, status, created_by, has_login
		) VALUES (
			@id, @staff_id, @password_hash, @first_name, @last_name, @email, @phone, @address, @date_of_birth, @date_of_hire, @emergency_contact_name, @emergency_contact_phone, @bank_name, @bank_account_number, @bank_account_name, @department, @job_title, @pay_type, @base_salary, @status, @created_by, @has_login
		) RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, created_at, updated_at, status
	`

	var exists bool
	checkStmt := `SELECT EXISTS(SELECT 1 FROM staff WHERE staff_id = $1)`
	err := s.db.QueryRow(ctx, checkStmt, credentials.StaffID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if staff exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("staff_id %s already exists", credentials.StaffID)
	}

	rows, err := s.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":                      id,
		"staff_id":                credentials.StaffID,
		"password_hash":           credentials.Password,
		"first_name":              payload.FirstName,
		"last_name":               payload.LastName,
		"email":                   &payload.Email,
		"phone":                   &payload.Phone,
		"address":                 &payload.Address,
		"date_of_birth":           &payload.DateOfBirth,
		"date_of_hire":            payload.DateOfHire,
		"emergency_contact_name":  &payload.EmergencyContactName,
		"emergency_contact_phone": &payload.EmergencyContactPhone,
		"bank_name":               &payload.BankName,
		"bank_account_number":     &payload.BankAccountNumber,
		"bank_account_name":       &payload.BankAccountName,
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

	staff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.StaffCreate])
	if err != nil {
		return nil, fmt.Errorf("failed to collect created staff: %w", err)
	}

	return &staff, nil
}

// GetStaffByID retrieves a staff member by ID
func (s *StaffRepository) GetStaffByID(ctx context.Context, id uuid.UUID) (*staff.Staff, error) {
	stmt := `
		SELECT
			s.id, s.staff_id, s.password_hash, s.first_name, s.last_name, s.email, s.phone, s.address, s.date_of_birth, s.date_of_hire,
			s.emergency_contact_name, s.emergency_contact_phone, s.bank_name, s.bank_account_number,
			s.bank_account_name, s.department, s.job_title, s.pay_type, s.base_salary, s.status, s.fired_at, s.has_login,
			s.supabase_uid, s.created_at, s.updated_at, creator.id AS creator_id, creator.staff_id AS creator_staff_id, creator.first_name AS creator_first_name, creator.last_name AS creator_last_name, creator.email AS creator_email, creator.phone AS creator_phone, creator.department AS creator_department, creator.job_title AS creator_job_title, creator.status AS creator_status
		FROM staff s
		LEFT JOIN staff creator ON s.created_by = creator.id
		WHERE s.id = $1
	`

	rows, err := s.db.Query(ctx, stmt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting staff by ID: %w", err)
	}

	staff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.Staff])
	if err != nil {
		return nil, fmt.Errorf("failed to collect staff by ID: %w", err)
	}

	return &staff, nil
}

func (s *StaffRepository) GetStaffByStaffID(ctx context.Context, staffID string) (*staff.Staff, error) {
	stmt := `
		SELECT
			s.id, s.staff_id, s.password_hash, s.first_name, s.last_name, s.email, s.phone, s.address, s.date_of_birth, s.date_of_hire,
			s.emergency_contact_name, s.emergency_contact_phone, s.bank_name, s.bank_account_number,
			s.bank_account_name, s.department, s.job_title, s.pay_type, s.base_salary, s.status, s.fired_at, s.has_login,
			s.supabase_uid, s.created_at, s.updated_at, creator.id AS creator_id, creator.staff_id AS creator_staff_id, creator.first_name AS creator_first_name, creator.last_name AS creator_last_name, creator.email AS creator_email, creator.phone AS creator_phone, creator.department AS creator_department, creator.job_title AS creator_job_title, creator.status AS creator_status
		FROM staff s
		LEFT JOIN staff creator ON s.created_by = creator.id
		WHERE s.staff_id = $1
	`

	rows, err := s.db.Query(ctx, stmt, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query for getting staff by staff_id: %w", err)
	}

	staff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.Staff])
	if err != nil {
		return nil, fmt.Errorf("failed to find staff by staff_id, not found: %w", err)
	}

	return &staff, nil
}

// GetAllStaff retrieves all staff members with pagination
func (s *StaffRepository) GetAllStaff(ctx context.Context, limit, offset int) ([]staff.Staff, error) {
	stmt := `
		SELECT
			s.id, s.password_hash, s.staff_id, s.first_name, s.last_name, s.email, s.phone, s.address, s.date_of_birth, s.date_of_hire,
			s.emergency_contact_name, s.emergency_contact_phone, s.bank_name, s.bank_account_number,
			s.bank_account_name, s.department, s.job_title, s.pay_type, s.base_salary, s.status, s.fired_at,
			s.has_login, s.supabase_uid, s.created_at, s.updated_at, creator.id AS creator_id, creator.staff_id AS creator_staff_id, creator.first_name AS creator_first_name, creator.last_name AS creator_last_name, creator.email AS creator_email, creator.phone AS creator_phone, creator.department AS creator_department, creator.job_title AS creator_job_title, creator.status AS creator_status
		FROM staff s
		LEFT JOIN staff creator ON s.created_by = creator.id
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, stmt, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all staff: %w", err)
	}

	staffList, err := pgx.CollectRows(rows, pgx.RowToStructByName[staff.Staff])
	if err != nil {
		return nil, fmt.Errorf("failed to collect staff: %w", err)
	}

	return staffList, nil
}

// UpdateStaff updates staff member information
func (s *StaffRepository) UpdateStaff(ctx context.Context, id uuid.UUID, payload *staff.UpdateStaffPayload) (*staff.StaffCreate, error) {
	stmt := `
		UPDATE staff
		SET
		    first_name = COALESCE(@first_name, first_name),
		    last_name = COALESCE(@last_name, last_name),
		    email = COALESCE(@email, email),
		    phone = COALESCE(@phone, phone),
		    address = COALESCE(@address, address),
		    date_of_birth = COALESCE(@date_of_birth, date_of_birth),
		    emergency_contact_name = COALESCE(@emergency_contact_name, emergency_contact_name),
		    emergency_contact_phone = COALESCE(@emergency_contact_phone, emergency_contact_phone),
		    bank_name = COALESCE(@bank_name, bank_name),
		    bank_account_number = COALESCE(@bank_account_number, bank_account_number),
		    bank_account_name = COALESCE(@bank_account_name, bank_account_name),
		    department = COALESCE(@department, department),
		    job_title = COALESCE(@job_title, job_title),
		    pay_type = COALESCE(@pay_type, pay_type),
		    base_salary = COALESCE(@base_salary, base_salary),
		    status = COALESCE(@status, status),
		    updated_at = NOW()
		WHERE id = @id
		RETURNING id, staff_id, first_name, last_name, email, phone, department, job_title, status, created_at, updated_at
	`

	rows, err := s.db.Query(ctx, stmt, pgx.NamedArgs{
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
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	updatedStaff, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[staff.StaffCreate])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated staff: %w", err)
	}

	return &updatedStaff, nil
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

	var roleList []roles.RoleSummary
	for rows.Next() {
		var role roles.RoleSummary
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roleList = append(roleList, role)
	}
	return roleList, nil
}

func (s *StaffRepository) GetStaffPermissions(ctx context.Context, id uuid.UUID) ([]roles.PermissionSummary, error) {
	stmt := `
		SELECT DISTINCT p.id, p.code, p.module
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN staff_roles sr ON rp.role_id = sr.role_id
		WHERE sr.staff_id = $1
		ORDER BY p.module, p.code
      `

	rows, err := s.db.Query(ctx, stmt, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch staff permissions: %w", err)
	}

	defer rows.Close()

	var permList []roles.PermissionSummary
	for rows.Next() {
		var perm roles.PermissionSummary
		err := rows.Scan(&perm.ID, &perm.Code, &perm.Module)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permList = append(permList, perm)
	}
	return permList, nil
}
