package main

import (
	"log"

	"github.com/bekbek22/JaiYenMarket_backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectMongo()

	app := fiber.New()

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

	app.Listen(":3000")
}
