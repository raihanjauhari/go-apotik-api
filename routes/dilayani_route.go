package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func DilayaniRoute(app *fiber.App) {
	dilayani := app.Group("/api/dilayani")

	dilayani.Get("/", handlers.GetAllDilayani)           // GET /dilayani - ambil semua data
	dilayani.Get("/:id", handlers.GetDilayaniByID)       // GET /dilayani/:id - ambil data by ID_PENDAFTARAN
	dilayani.Post("/", handlers.CreateDilayani)           // POST /dilayani - tambah data baru
	dilayani.Put("/:id", handlers.UpdateDilayani)        // PUT /dilayani/:id - update data by ID_PENDAFTARAN
	dilayani.Delete("/:id", handlers.DeleteDilayani)     // DELETE /dilayani/:id - hapus data by ID_PENDAFTARAN
}
