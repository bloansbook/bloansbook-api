package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/bloansbook/bloansbook-api/internal/suppliers/handler"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router, h *handler.SupplierHandler, ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) {
	suppliers := api.Group("/suppliers", middleware.Auth(ar, sr))

	suppliers.Get("/", middleware.RequirePermission("suppliers.view"), h.GetAllSuppliers)
	suppliers.Post("/", middleware.RequirePermission("suppliers.create"), h.CreateSupplier)
	suppliers.Get("/:id", middleware.RequirePermission("suppliers.view"), h.GetSupplierById)
	suppliers.Patch("/:id", middleware.RequirePermission("suppliers.update"), h.UpdateSupplier)
}
