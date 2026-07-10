package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/handler"
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router, h *handler.AuthHandler, ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) {
	auth := api.Group("/auth")

	// Public — no middleware
	auth.Post("/login", h.Login)

	// Protected — token required for everything below
	protected := auth.Group("", middleware.Auth(ar, sr))

	// Only super_admin can manage auth accounts and reset passwords
	protected.Post("/accounts", middleware.RequirePermission("auth.manage_accounts"), h.Login) // placeholder — swap for real handler
	protected.Post("/reset-password", middleware.RequirePermission("auth.reset_password"), h.Login)
	protected.Post("/manage-roles", middleware.RequirePermission("auth.manage_roles"), h.Login)
}
