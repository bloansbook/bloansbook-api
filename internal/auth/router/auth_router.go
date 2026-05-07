package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/handler"
	"github.com/gofiber/fiber/v3"
)

// SetupRolesRoutes registers all roles-related routes
func SetupRoutes(api fiber.Router, h *handler.AuthHandler) {
	staff := api.Group("/auth")

	// Auth management
	staff.Post("/login", h.Login)
}
