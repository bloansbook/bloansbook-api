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

	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)

	protected := auth.Group("", middleware.Auth(ar, sr))
	protected.Get("/me", h.GetProfile)
}
