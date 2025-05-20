package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func DetailEresepRoute(app *fiber.App) {
    detailEresep := app.Group("/api/detail_eresep")

    detailEresep.Get("/", handlers.GetAllDetailEresep)       // GET all
    detailEresep.Get("/:id", handlers.GetDetailEresepByID)    // GET by ID
    detailEresep.Post("/", handlers.CreateDetailEresep)       // POST create
    detailEresep.Put("/:id", handlers.UpdateDetailEresep)     // PUT update by ID
    detailEresep.Delete("/:id", handlers.DeleteDetailEresep)  // DELETE by ID
}
