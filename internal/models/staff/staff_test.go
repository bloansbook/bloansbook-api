package staff_test

import (
	"testing"
	"time"

	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func sampleStaff() staff.Staff {
	now := time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC)
	email := "amara@example.com"
	creatorEmail := "system@example.com"

	return staff.Staff{
		BaseModel: models.BaseModel{
			BaseWithId:        models.BaseWithId{ID: uuid.New()},
			BaseWithCreatedAt: models.BaseWithCreatedAt{CreatedAt: now},
			BaseWithUpdatedAt: models.BaseWithUpdatedAt{UpdatedAt: now},
		},
		StaffID:           "BLN-0001",
		FirstName:         "Amara",
		LastName:          "Okafor",
		Email:             &email,
		Department:        "management",
		JobTitle:          "General Manager",
		PayType:           "monthly",
		BaseSalary:        decimal.NewFromInt(350000),
		Status:            staff.StaffStatusActive,
		CreatorStaffID:    "SYSTEM-0000",
		CreatorFirstName:  "System",
		CreatorLastName:   "Account",
		CreatorEmail:      &creatorEmail,
		CreatorDepartment: "system",
		CreatorJobTitle:   "system_account",
		CreatorStatus:     staff.StaffStatusActive,
	}
}

func TestStaff_ToSummary(t *testing.T) {
	s := sampleStaff()
	sum := s.ToSummary()

	if sum.StaffID != s.StaffID || sum.FirstName != s.FirstName || sum.LastName != s.LastName {
		t.Errorf("identity fields not mapped: %+v", sum)
	}
	if sum.Department != s.Department || sum.JobTitle != s.JobTitle {
		t.Errorf("department/jobTitle not mapped: %+v", sum)
	}
	if sum.Status != s.Status {
		t.Errorf("Status: got %q want %q", sum.Status, s.Status)
	}
}

func TestStaff_ToDTO_NilRolesBecomesEmptySlice(t *testing.T) {
	s := sampleStaff()
	dto := s.ToDTO(nil)

	if dto.Roles == nil {
		t.Fatal("expected non-nil empty Roles slice, got nil")
	}
	if len(dto.Roles) != 0 {
		t.Errorf("expected empty Roles, got %d", len(dto.Roles))
	}
	if dto.CreatedBy.StaffID != s.CreatorStaffID {
		t.Errorf("CreatedBy.StaffID: got %q want %q", dto.CreatedBy.StaffID, s.CreatorStaffID)
	}
	if dto.BaseSalary.Cmp(s.BaseSalary) != 0 {
		t.Errorf("BaseSalary: got %s want %s", dto.BaseSalary, s.BaseSalary)
	}
}

func TestStaff_ToDTO_PassesThroughRoles(t *testing.T) {
	s := sampleStaff()
	dto := s.ToDTO([]roles.RoleSummary{{}, {}})

	if len(dto.Roles) != 2 {
		t.Errorf("expected 2 roles passed through, got %d", len(dto.Roles))
	}
}
