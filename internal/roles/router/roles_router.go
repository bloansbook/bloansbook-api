package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	"github.com/bloansbook/bloansbook-api/internal/roles/handler"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router, h *handler.RolesHandler, ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) {
	roles := api.Group("/roles", middleware.Auth(ar, sr))

	// Role management — requires auth.manage_roles
	roles.Post("/", middleware.RequirePermission("auth.manage_roles"), h.CreateRole)
	roles.Get("/", middleware.RequirePermission("auth.manage_roles"), h.GetAllRoles)
	roles.Get("/:id", middleware.RequirePermission("auth.manage_roles"), h.GetRoleWithPermissions)

	// Permission management — requires auth.manage_roles
	roles.Post("/permissions", middleware.RequirePermission("auth.manage_roles"), h.CreatePermission)

	// Assign permissions to roles — requires auth.manage_roles
	roles.Post("/assign-permission", middleware.RequirePermission("auth.manage_roles"), h.AssignPermissionToRole)
}
