package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func MemunculkanRoute(app *fiber.App) {
    memunculkan := app.Group("/api/memunculkan")

    memunculkan.Get("/", handlers.GetAllMemunculkan)
    memunculkan.Get("/:kode_obat/:id_eresep", handlers.GetMemunculkanByIDs)
    memunculkan.Post("/", handlers.CreateMemunculkan)
    memunculkan.Put("/:kode_obat/:id_eresep", handlers.UpdateMemunculkan)
    memunculkan.Delete("/:kode_obat/:id_eresep", handlers.DeleteMemunculkan)
}
