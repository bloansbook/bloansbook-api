package router

import (
	"github.com/bloansbook/bloansbook-api/internal/staff/handler"
	"github.com/gofiber/fiber/v3"
)

// SetupRolesRoutes registers all roles-related routes
func SetupRoutes(api fiber.Router, h *handler.StaffHandler) {
	staff := api.Group("/staff")

	// Staff management
	staff.Get("/", h.GetAllStaff)
	staff.Post("/", h.CreateStaff)
	staff.Get("/:id", h.GetStaffById)
	staff.Patch("/:id", h.UpdateStaff)
}
