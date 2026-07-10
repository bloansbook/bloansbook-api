package middleware

import (
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
)

// RequirePermission returns a handler that blocks the request if the
// authenticated staff member does not hold the given permission code.
//
// Must be used on routes that are already protected by Auth middleware,
// since it reads from ctx.Locals populated there.
//
// Usage:
//
//	staff.Post("/", middleware.RequirePermission("staff.create"), h.CreateStaff)
//	roles.Delete("/:id", middleware.RequirePermission("auth.manage_roles"), h.DeleteRole)
func RequirePermission(code string) fiber.Handler {
	return func(c fiber.Ctx) error {
		perms, ok := c.Locals(LocalPermissions).([]string)
		if !ok {
			return response.Error(c, sysmsg.Forbidden, fiber.StatusForbidden)
		}

		for _, p := range perms {
			if p == code {
				return c.Next()
			}
		}

		return response.Error(c, sysmsg.Forbidden, fiber.StatusForbidden)
	}
}

// RequireRole blocks the request if the authenticated staff member does not
// hold the given role name. Use this for coarse-grained checks
// (e.g. super_admin only routes) where a permission code is not warranted.
//
// Usage:
//
//	settings.Patch("/", middleware.RequireRole("super_admin"), h.UpdateSetting)
func RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roles, ok := c.Locals(LocalRoles).([]string)
		if !ok {
			return response.Error(c, sysmsg.Forbidden, fiber.StatusForbidden)
		}

		for _, r := range roles {
			if r == role {
				return c.Next()
			}
		}

		return response.Error(c, sysmsg.Forbidden, fiber.StatusForbidden)
	}
}

// CallerStaffID returns the authenticated staff member's UUID from locals.
// Returns the zero value if the middleware has not run.
func CallerStaffID(c fiber.Ctx) (interface{}, bool) {
	id := c.Locals(LocalStaffID)
	return id, id != nil
}

// CallerStaffIDStr returns the human-readable staff ID (e.g. "BLN-0001").
func CallerStaffIDStr(c fiber.Ctx) string {
	s, _ := c.Locals(LocalStaffIDStr).(string)
	return s
}

// HasPermission reports whether the caller holds a specific permission.
// Useful for conditional logic inside handlers without blocking the whole request.
func HasPermission(c fiber.Ctx, code string) bool {
	perms, ok := c.Locals(LocalPermissions).([]string)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == code {
			return true
		}
	}
	return false
}
