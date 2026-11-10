package main

import (
	"keuangan_app/database"
	"keuangan_app/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()

	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))

}