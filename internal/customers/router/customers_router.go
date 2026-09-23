package router

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	"github.com/bloansbook/bloansbook-api/internal/customers/handler"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router, h *handler.CustomerHandler, ar *authRepo.AuthRepository, sr *staffRepo.StaffRepository) {
	customers := api.Group("/customers", middleware.Auth(ar, sr))

	customers.Get("/", middleware.RequirePermission("customers.view"), h.GetAllCustomers)
	customers.Post("/", middleware.RequirePermission("customers.create"), h.CreateCustomer)
	customers.Get("/:id", middleware.RequirePermission("customers.view"), h.GetCustomerById)
	customers.Patch("/:id", middleware.RequirePermission("customers.update"), h.UpdateCustomer)
}
