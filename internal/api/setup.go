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

func SetupRoutes(app *fiber.App) {
	db := database.Pool
	cfg := config.ApplicationConfig

	api := app.Group("/api/v1")

	authRepository := authRepo.NewAuthRepository(cfg, db)
	staffRepository := staffRepo.NewStaffRepository(db, cfg)

	staffUC := staffUsecase.NewStaffUsecase(db, authRepository, staffRepository, cfg)
	authUC := authUsecase.NewAuthUsecase(authRepository, staffRepository, staffUC)
	authRouter.SetupRoutes(api, authHandler.NewAuthHandler(authUC), authRepository, staffRepository)

	rolesRepository := rolesRepo.NewRolesRepository(db, cfg)
	rolesUC := rolesUsecase.NewRolesUsecase(rolesRepository)
	rolesRouter.SetupRoutes(api, rolesHandler.NewRolesHandler(rolesUC), authRepository, staffRepository)

	staffRouter.SetupRoutes(api, staffHandler.NewStaffHandler(staffUC), authRepository, staffRepository)
}
