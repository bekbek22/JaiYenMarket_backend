package routes

import (
	"github.com/bekbek22/JaiYenMarket_backend/handler"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoutes sets up routes for authentication
func RegisterAuthRoutes(app *fiber.App, authHandler handler.IAuthHandler) {
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/register", authHandler.Register) // POST /auth/register
	authGroup.Post("/login", authHandler.Login)       // POST /auth/login
}

// RegisterTradeRoutes sets up routes for trade
func RegisterTradeRoutes(app *fiber.App, tradeHandler handler.ITradeHandler) {
	tradeGroup := app.Group("/api/v1/trade", middleware.JWTMiddleware("supersecret"))
	tradeGroup.Post("/", tradeHandler.CreateTrade)
	tradeGroup.Get("/me", tradeHandler.GetMyTrades)
	tradeGroup.Post("/:id/add-item", tradeHandler.AddItemToTrade)
	tradeGroup.Post("/:id/confirm", tradeHandler.ConfirmTrade)
	tradeGroup.Post("/:id/unconfirm", tradeHandler.UnconfirmTrade)
	tradeGroup.Post("/:id/cancel", tradeHandler.CancelTrade)
	tradeGroup.Get("/:id", tradeHandler.GetTrade)
}

func RegisterWalletRoutes(app *fiber.App, walletHandler handler.IWalletHandler) {
	w := app.Group("/api/v1/wallet", middleware.JWTMiddleware("supersecret"))
	w.Get("/", walletHandler.GetBalance)
	w.Post("/deposit", walletHandler.Deposit)
	w.Post("/withdraw", walletHandler.Withdraw)
	w.Post("/transfer", walletHandler.Transfer)
}

// Register all routes sets up for application
func RegisterRoutes(app *fiber.App, authHandler handler.IAuthHandler, tradeHandler handler.ITradeHandler, walletHandler handler.IWalletHandler) {
	RegisterAuthRoutes(app, authHandler)
	RegisterTradeRoutes(app, tradeHandler)
	RegisterWalletRoutes(app, walletHandler)
}
