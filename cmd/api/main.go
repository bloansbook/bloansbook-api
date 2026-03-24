package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/database"
)

func main() {
	// Load config
	config.Load()

	// Connect to DB
	database.Connect()
	defer database.Close()

	app := fiber.New(fiber.Config{
		AppName: "BloansBook API v1",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "BloansBooks API",
		})
	})

	api := app.Group("/api/v1")
	_ = api

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		port := config.ApplicationConfig.App.Port
		if port == "" {
			port = "8080"
		}
		log.Printf("BloansBooks API running on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	_ = app.Shutdown()
	log.Println("Server stopped")
}
