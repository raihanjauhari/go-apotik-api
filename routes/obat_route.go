package route

import (
	"database/sql"
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterObatRoutes(app *fiber.App, db *sql.DB) {
	obat := app.Group("/api/obat") // base path api/obat

	obat.Get("/", handlers.GetAllObat(db))          
	obat.Get("/:kode", handlers.GetObatByKode(db))   
	obat.Post("/", handlers.CreateObat(db))          
	obat.Put("/:kode", handlers.UpdateObat(db))      
	obat.Delete("/:kode", handlers.DeleteObat(db))   
}
