package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/bloansbook/bloansbook-api/internal/api"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/database"
)

func main() {
	config.Load()

	database.Connect()
	defer database.Close()

	app := fiber.New(fiber.Config{
		AppName: "BloansBook API v1",
	})

	rawOrigins := config.ApplicationConfig.App.AllowedOrigins

	var allowedOrigins []string
	allowCredentials := false

	if rawOrigins == "" || rawOrigins == "*" {
		allowedOrigins = []string{"*"}
	} else {
		for _, o := range strings.Split(rawOrigins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
		allowCredentials = true
	}

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: allowCredentials,
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "BloansBooks API",
		})
	})

	api.SetupRoutes(app)

	serverErrors := make(chan error, 1)

	go func() {
		port := config.ApplicationConfig.App.Port
		if port == "" {
			port = "8080"
		}
		log.Printf("BloansBooks API running on port %s", port)
		serverErrors <- app.Listen(":" + port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting server: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v. Shutting down server...", sig)
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
	}

	log.Println("Server gracefully stopped")
}
