package api

import (
	authHandler "github.com/bloansbook/bloansbook-api/internal/auth/handler"
	authRepo "github.com/bloansbook/bloansbook-api/internal/auth/repository"
	authRouter "github.com/bloansbook/bloansbook-api/internal/auth/router"
	authUsecase "github.com/bloansbook/bloansbook-api/internal/auth/usecase"

	rolesHandler "github.com/bloansbook/bloansbook-api/internal/roles/handler"
	rolesRepo "github.com/bloansbook/bloansbook-api/internal/roles/repository"
	rolesRouter "github.com/bloansbook/bloansbook-api/internal/roles/router"
	rolesUsecase "github.com/bloansbook/bloansbook-api/internal/roles/usecase"

	staffHandler "github.com/bloansbook/bloansbook-api/internal/staff/handler"
	staffRepo "github.com/bloansbook/bloansbook-api/internal/staff/repository"
	staffRouter "github.com/bloansbook/bloansbook-api/internal/staff/router"
	staffUsecase "github.com/bloansbook/bloansbook-api/internal/staff/usecase"

	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/database"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes initializes all API route modules
func SetupRoutes(app *fiber.App) {
	// Get shared dependencies
	db := database.Pool
	cfg := config.ApplicationConfig

	// Create API group
	api := app.Group("/api/v1")

	// ===== AUTH MODULE (TODO: Add after auth implementation) =====
	authRepository := authRepo.NewAuthRepository(cfg)
	staffRepo := staffRepo.NewStaffRepository(db, cfg)
	staffUsecase := staffUsecase.NewStaffUsecase(db, authRepository, staffRepo, cfg)
	authUsecase := authUsecase.NewAuthUsecase(authRepository, staffRepo, staffUsecase)
	authHandler := authHandler.NewAuthHandler(authUsecase)
	authRouter.SetupRoutes(api, authHandler)

	// ===== ROLES MODULE =====
	rolesRepo := rolesRepo.NewRolesRepository(db, cfg)
	rolesUsecase := rolesUsecase.NewRolesUsecase(rolesRepo)
	rolesHandler := rolesHandler.NewRolesHandler(rolesUsecase)
	rolesRouter.SetupRoutes(api, rolesHandler)

	// ===== STAFF MODULE (TODO: Add after staff implementation) =====
	staffHandler := staffHandler.NewStaffHandler(staffUsecase)
	staffRouter.SetupRoutes(api, staffHandler)

	// ===== OTHER MODULES =====
	// Add more modules here as you build them:
	// - inventory
	// - sales
	// - customers
	// - products
	// - invoices
	// - payments
	// - suppliers
	// - purchase_orders
	// - bills
	// - materials
	// - job_costing
	// - payroll
	// - attendance
	// - tasks
	// - notifications
	// - reports
	// - dashboard
	// - audit
}
