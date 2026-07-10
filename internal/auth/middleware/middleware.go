package middleware

import (
	"strings"

	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
)

// Context keys stored in fiber.Ctx.Locals.
const (
	LocalStaffID    = "staff_id"    // uuid.UUID
	LocalStaffIDStr = "staff_id_str" // string  e.g. "BLN-0001"
	LocalRoles      = "roles"       // []string
	LocalPermissions = "permissions" // []string  e.g. ["staff.create", "invoices.post"]
)

// Auth validates the Bearer JWT issued by Supabase, resolves the staff record,
// and loads their roles and permissions into ctx.Locals.
//
// Usage — add to any route group that requires authentication:
//
//	protected := api.Group("/staff", middleware.Auth(authRepository, staffRepository))
func Auth(ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Extract Bearer token from Authorization header
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
		}
		token := strings.TrimPrefix(header, "Bearer ")

		// Verify token with Supabase and get the caller's identity
		user, err := ar.GetUser(c.Context(), token)
		if err != nil {
			return response.Error(c, sysmsg.TokenInvalid, fiber.StatusUnauthorized)
		}

		supabaseUID := user.ID.String()

		// Resolve staff record by supabase_uid
		staff, err := sr.GetStaffBySupabaseUID(c.Context(), supabaseUID)
		if err != nil {
			return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
		}

		// Block fired / inactive staff
		if staff.Status != "active" {
			return response.Error(c, sysmsg.AccessDenied, fiber.StatusForbidden)
		}

		// Load this staff member's permissions
		permissions, err := sr.GetStaffPermissions(c.Context(), staff.ID)
		if err != nil {
			// Non-fatal — staff may simply have no permissions yet
			permissions = []string{}
		}

		// Load role names
		roles, err := sr.GetStaffRoles(c.Context(), staff.ID)
		if err != nil {
			roles = nil
		}
		roleNames := make([]string, len(roles))
		for i, r := range roles {
			roleNames[i] = r.Name
		}

		// Store identity in context locals for handlers and guards
		c.Locals(LocalStaffID, staff.ID)
		c.Locals(LocalStaffIDStr, staff.StaffID)
		c.Locals(LocalRoles, roleNames)
		c.Locals(LocalPermissions, permissions)

		return c.Next()
	}
}
