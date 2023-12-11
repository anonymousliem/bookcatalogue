package main

import (
	"book-catalogue/configs"
	"book-catalogue/routes"
	"book-catalogue/controllers"
	"log"


	"github.com/gofiber/fiber/v2"
)


func main() {
	app := fiber.New()
	configs.ConnectDB()
	controllers.InitializeRedisClient()
	routes.BookRoute(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "it works"})
	})

	// Listen on port 6000
	err := app.Listen(":6000")
	if err != nil {
		log.Fatal(err)
	}
}
