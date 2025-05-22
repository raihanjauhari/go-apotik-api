package main

import (
	"go-apotik-api/database"
	route "go-apotik-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Tambahkan middleware CORS di sini
	app.Use(cors.New(cors.Config{
	AllowOrigins:     "http://localhost:5174",
	AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	AllowHeaders:     "*", // atau bisa sebutkan secara spesifik
	AllowCredentials: true,
}))


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
