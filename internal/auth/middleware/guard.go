package middleware

import (
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
)

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

func CallerStaffID(c fiber.Ctx) (interface{}, bool) {
	id := c.Locals(LocalStaffID)
	return id, id != nil
}

func CallerStaffIDStr(c fiber.Ctx) string {
	s, _ := c.Locals(LocalStaffIDStr).(string)
	return s
}

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
