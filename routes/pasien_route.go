package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func PasienRoute(app *fiber.App) {
    r := app.Group("/api/pasien")
    r.Get("/", handlers.GetAllPasien)
    r.Get("/:id", handlers.GetPasienByID)
    r.Post("/", handlers.CreatePasien)
    r.Put("/:id", handlers.UpdatePasien)
    r.Delete("/:id", handlers.DeletePasien)
}
