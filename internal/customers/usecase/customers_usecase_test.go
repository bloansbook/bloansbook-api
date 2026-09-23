package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/bloansbook/bloansbook-api/internal/models/customers"
	"github.com/google/uuid"
)

// The validation guards in CreateCustomer/UpdateCustomer run before any
// repository or database access, so they can be exercised with a zero-value
// usecase (nil repo, nil db). Every case here must trip a guard; an input that
// passes validation would fall through to the DB and belongs in an integration
// test, not here.

func TestCreateCustomer_GuardRejections(t *testing.T) {
	validCaller := uuid.New()

	tests := []struct {
		name    string
		caller  uuid.UUID
		payload customers.CreateCustomerPayload
		wantErr string
	}{
		{
			name:    "missing caller",
			caller:  uuid.Nil,
			payload: customers.CreateCustomerPayload{Name: "Ada", Phone: "+2348012345678", Type: customers.CustomerTypeRetail},
			wantErr: "authentication required",
		},
		{
			name:    "empty name",
			caller:  validCaller,
			payload: customers.CreateCustomerPayload{Name: "", Phone: "+2348012345678", Type: customers.CustomerTypeRetail},
			wantErr: "name is required",
		},
		{
			name:    "whitespace-only name",
			caller:  validCaller,
			payload: customers.CreateCustomerPayload{Name: "   ", Phone: "+2348012345678", Type: customers.CustomerTypeRetail},
			wantErr: "name is required",
		},
		{
			name:    "empty phone",
			caller:  validCaller,
			payload: customers.CreateCustomerPayload{Name: "Ada", Phone: "  ", Type: customers.CustomerTypeRetail},
			wantErr: "phone is required",
		},
		{
			name:    "empty type",
			caller:  validCaller,
			payload: customers.CreateCustomerPayload{Name: "Ada", Phone: "+2348012345678", Type: ""},
			wantErr: "type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CustomerUsecase{}
			payload := tt.payload
			_, err := u.CreateCustomer(context.Background(), tt.caller, &payload)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// TestCreateCustomer_TrimsRequiredFields proves the trim mutates the payload:
// a padded-but-valid name passes its guard, and the next guard (phone) fires,
// leaving the trimmed name observable on the payload.
func TestCreateCustomer_TrimsRequiredFields(t *testing.T) {
	u := &CustomerUsecase{}
	payload := customers.CreateCustomerPayload{
		Name:  "  Ada Lovelace  ",
		Phone: "   ",
		Type:  customers.CustomerTypeRetail,
	}

	_, err := u.CreateCustomer(context.Background(), uuid.New(), &payload)
	if err == nil || !strings.Contains(err.Error(), "phone is required") {
		t.Fatalf("expected phone guard to fire after name trim, got %v", err)
	}
	if payload.Name != "Ada Lovelace" {
		t.Fatalf("expected name trimmed to %q, got %q", "Ada Lovelace", payload.Name)
	}
}

func TestUpdateCustomer_GuardRejections(t *testing.T) {
	blank := "   "
	empty := ""
	blankType := customers.CustomerType("  ")

	tests := []struct {
		name    string
		payload customers.UpdateCustomerPayload
		wantErr string
	}{
		{
			name:    "blank name",
			payload: customers.UpdateCustomerPayload{Name: &blank},
			wantErr: "name cannot be empty",
		},
		{
			name:    "empty phone",
			payload: customers.UpdateCustomerPayload{Phone: &empty},
			wantErr: "phone cannot be empty",
		},
		{
			name:    "blank type",
			payload: customers.UpdateCustomerPayload{Type: &blankType},
			wantErr: "type cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &CustomerUsecase{}
			payload := tt.payload
			_, err := u.UpdateCustomer(context.Background(), uuid.New(), &payload)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
