package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kunalsinghdadhwal/fib_notes/handlers"
	"github.com/kunalsinghdadhwal/fib_notes/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	authHandler := handlers.NewAuthHandler()

	auth := app.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	auth.Post("/logout", middleware.JWTMiddleware(), authHandler.Logout)
	auth.Get("/me", middleware.JWTMiddleware(), authHandler.Me)
	auth.Put("/change-password", middleware.JWTMiddleware(), authHandler.ChangePassword)
}
