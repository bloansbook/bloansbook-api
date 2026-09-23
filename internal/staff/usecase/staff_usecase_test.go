package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/google/uuid"
)

// CreateStaff/UpdateStaff trim and validate required string fields before any
// repository or database call, so a zero-value usecase suffices. Each case must
// trip a guard; a fully valid payload would fall through to the DB.

func validCreateStaffPayload() staff.CreateStaffPayload {
	return staff.CreateStaffPayload{
		FirstName:  "Ada",
		LastName:   "Lovelace",
		Department: "sales",
		JobTitle:   "Sales Executive",
		PayType:    "monthly",
	}
}

func TestCreateStaff_GuardRejections(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(p *staff.CreateStaffPayload)
		wantErr string
	}{
		{"empty firstName", func(p *staff.CreateStaffPayload) { p.FirstName = "  " }, "firstName is required"},
		{"empty lastName", func(p *staff.CreateStaffPayload) { p.LastName = "" }, "lastName is required"},
		{"empty department", func(p *staff.CreateStaffPayload) { p.Department = "   " }, "department is required"},
		{"empty jobTitle", func(p *staff.CreateStaffPayload) { p.JobTitle = "" }, "jobTitle is required"},
		{"empty payType", func(p *staff.CreateStaffPayload) { p.PayType = "  " }, "payType is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &StaffUsecase{}
			p := validCreateStaffPayload()
			tt.mutate(&p)
			_, err := u.CreateStaff(context.Background(), uuid.New(), &p)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestUpdateStaff_GuardRejections(t *testing.T) {
	ptr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		payload staff.UpdateStaffPayload
		wantErr string
	}{
		{"blank firstName", staff.UpdateStaffPayload{FirstName: ptr("  ")}, "firstName cannot be empty"},
		{"empty lastName", staff.UpdateStaffPayload{LastName: ptr("")}, "lastName cannot be empty"},
		{"blank department", staff.UpdateStaffPayload{Department: ptr("   ")}, "department cannot be empty"},
		{"empty jobTitle", staff.UpdateStaffPayload{JobTitle: ptr("")}, "jobTitle cannot be empty"},
		{"blank payType", staff.UpdateStaffPayload{PayType: ptr(" ")}, "payType cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &StaffUsecase{}
			payload := tt.payload
			_, err := u.UpdateStaff(context.Background(), uuid.New(), &payload)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
