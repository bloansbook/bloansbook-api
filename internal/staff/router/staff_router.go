package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	"github.com/bloansbook/bloansbook-api/internal/staff/handler"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router, h *handler.StaffHandler, ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) {
	staff := api.Group("/staff", middleware.Auth(ar, sr))

	staff.Get("/", middleware.RequirePermission("staff.view"), h.GetAllStaff)
	staff.Post("/", middleware.RequirePermission("staff.create"), h.CreateStaff)
	staff.Get("/:id", middleware.RequirePermission("staff.view"), h.GetStaffById)
	staff.Patch("/:id", middleware.RequirePermission("staff.update"), h.UpdateStaff)
	staff.Post("/:id/fire", middleware.RequirePermission("staff.terminate"), h.FireStaff)
	staff.Patch("/:id/fire/override", middleware.RequirePermission("staff.terminate"), h.OverrideTermination)

	staff.Post("/:id/roles", middleware.RequirePermission("auth.manage_roles"), h.AssignRole)
	staff.Delete("/:id/roles", middleware.RequirePermission("auth.manage_roles"), h.RevokeRole)
	staff.Put("/:id/roles", middleware.RequirePermission("auth.manage_roles"), h.UpdateRole)
	staff.Get("/:id/roles/history", middleware.RequirePermission("auth.manage_roles"), h.GetRoleHistory)
}
