package main

import (
	"log"

	"github.com/bekbek22/JaiYenMarket_backend/config"
	"github.com/bekbek22/JaiYenMarket_backend/db"
	"github.com/bekbek22/JaiYenMarket_backend/handler"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/repository"
	"github.com/bekbek22/JaiYenMarket_backend/pkg/service"
	"github.com/bekbek22/JaiYenMarket_backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Load Config
	cfg := config.Load()

	db.ConnectMongo()

	app := fiber.New()

	userCol := db.MongoClient.Database("JaiYenMarket").Collection("users")
	tradeCol := db.MongoClient.Database("JaiYenMarket").Collection("trades")
	walletCol := db.MongoClient.Database("JaiYenMarket").Collection("wallets")
	transactionCol := db.MongoClient.Database("JaiYenMarket").Collection("transactions")

	authRepo := repository.NewAuthRepository(userCol)    // type: IAuthRepository
	authService := service.NewAuthService(authRepo, cfg) // type: IAuthService
	authHandler := handler.NewAuthHandler(authService)   // type: IAuthHandler

	tradeRepo := repository.NewTradeRepository(tradeCol)
	tradeService := service.NewTradeService(tradeRepo)
	tradeHandler := handler.NewTradeHandler(tradeService)

	walletRepo := repository.NewWalletRepository(walletCol, transactionCol)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	routes.RegisterRoutes(app, authHandler, tradeHandler, walletHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		collection := db.MongoClient.Database("mydb").Collection("users")
		return c.JSON(fiber.Map{
			"message":    "Mongo connected",
			"collection": collection.Name(),
		})
	})

	// app.Post("/register", authHandler.Register)
	// app.Post("/login", authHandler.Login)

	app.Listen(":3000")
}
