package middleware

import (
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
)

const (
	LocalStaffID     = "staff_id"
	LocalStaffIDStr  = "staff_id_str"
	LocalRoles       = "roles"
	LocalPermissions = "permissions"
)

func Auth(ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := GetAccessToken(c)
		if token == "" {
			return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
		}

		user, err := ar.GetUser(c.Context(), token)
		if err != nil {
			return response.Error(c, sysmsg.TokenInvalid, fiber.StatusUnauthorized)
		}

		supabaseUID := user.ID.String()

		staff, err := sr.GetStaffBySupabaseUID(c.Context(), supabaseUID)
		if err != nil {
			return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
		}

		if staff.Status != "active" {
			return response.Error(c, sysmsg.AccessDenied, fiber.StatusForbidden)
		}

		permissions, err := sr.GetStaffPermissions(c.Context(), staff.ID)
		if err != nil {
			permissions = []string{}
		}

		roles, err := sr.GetStaffRoles(c.Context(), staff.ID)
		if err != nil {
			roles = nil
		}
		roleNames := make([]string, len(roles))
		for i, r := range roles {
			roleNames[i] = r.Name
		}

		c.Locals(LocalStaffID, staff.ID)
		c.Locals(LocalStaffIDStr, staff.StaffID)
		c.Locals(LocalRoles, roleNames)
		c.Locals(LocalPermissions, permissions)

		return c.Next()
	}
}
