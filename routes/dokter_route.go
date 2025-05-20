package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func DokterRoute(app *fiber.App) {
	r := app.Group("/api/dokter") // ubah prefix route di sini
	r.Get("/", handlers.GetAllDokter)
	r.Get("/:id", handlers.GetDokterByID)
	r.Post("/", handlers.CreateDokter)
	r.Put("/:id", handlers.UpdateDokter)
	r.Delete("/:id", handlers.DeleteDokter)
}
