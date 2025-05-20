package main

import (
	"go-apotik-api/database"
	route "go-apotik-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	

	database.Connect()

	// Routes
	route.DokterRoute(app)
	route.PasienRoute(app)
	route.UserRoute(app)
	route.RegisterObatRoutes(app, database.DB)
	route.EresepRoute(app)
	route.DetailEresepRoute(app)

	app.Listen(":3000")
}
