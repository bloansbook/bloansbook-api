package router

import (
	"github.com/bloansbook/bloansbook-api/internal/staff/handler"
	"github.com/gofiber/fiber/v3"
)

// SetupRolesRoutes registers all roles-related routes
func SetupRoutes(api fiber.Router, h *handler.StaffHandler) {
	staff := api.Group("/staff")

	// Role management
	staff.Post("/", h.CreateStaff)
	staff.Get("/:id", h.GetStaffById)
	// roles.Get("/:id", h.GetRoleWithPermissions)

	// // Permission management
	// roles.Post("/permissions", h.CreatePermission)

	// // Assign permissions to roles
	// roles.Post("/assign-permission", h.AssignPermissionToRole)
}
