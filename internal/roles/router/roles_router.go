package router

import (
	"github.com/bloansbook/bloansbook-api/internal/roles/handler"
	"github.com/gofiber/fiber/v3"
)

// SetupRolesRoutes registers all roles-related routes
func SetupRoutes(api fiber.Router, h *handler.RolesHandler) {
	roles := api.Group("/roles")

	// Role management
	roles.Post("/", h.CreateRole)
	roles.Get("/", h.GetAllRoles)
	roles.Get("/:id", h.GetRoleWithPermissions)

	// Permission management
	roles.Post("/permissions", h.CreatePermission)

	// Assign permissions to roles
	roles.Post("/assign-permission", h.AssignPermissionToRole)
}
