package main

import (
	"go-apotik-api/database"
	route "go-apotik-api/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load file .env:", err)
	}

	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "*",
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
	route.DilayaniRoute(app)
	route.MemunculkanRoute(app)

	app.Static("/images", "./public/images")
	app.Listen(":3000")
}
