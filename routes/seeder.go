package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/fib_notes/handlers"
	"github.com/kunalsinghdadhwal/fib_notes/middleware"
)

func SetupSeederRoutes(app *fiber.App) {
	authHandler := handlers.NewAuthHandler()

	seeder := app.Group("/seeder")

	// Apply JWT middleware and Admin middleware to all seeder routes
	seeder.Use(middleware.JWTMiddleware())
	seeder.Use(middleware.AdminMiddleware())

	// Protected seeder endpoints (admin only)
	seeder.Post("/seed", authHandler.SeedDatabase)
	seeder.Post("/clear", authHandler.ClearDatabase)
	seeder.Get("/stats", authHandler.GetSeederStats)
}
