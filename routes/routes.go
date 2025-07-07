package routes

import (
	"github.com/bekbek22/JaiYenMarket_backend/handler"
	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoutes sets up routes for authentication
func RegisterAuthRoutes(app *fiber.App, authHandler handler.IAuthHandler) {
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/register", authHandler.Register) // POST /auth/register
	authGroup.Post("/login", authHandler.Login)       // POST /auth/login
}

// Register all routes sets up for application
func RegisterRoutes(app *fiber.App, authHandler handler.IAuthHandler) {
	RegisterAuthRoutes(app, authHandler)
}
