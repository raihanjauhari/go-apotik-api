package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func EresepRoute(app *fiber.App) {
    // Grouping route agar lebih rapi, misal prefix /eresep
    eresepRoutes := app.Group("/api/eresep")

    eresepRoutes.Post("/", handlers.CreateEresep)         // Create new eresep
    eresepRoutes.Get("/", handlers.GetAllEresep)           // Get all ereseps
    eresepRoutes.Get("/:id", handlers.GetEresepByID)       // Get eresep by ID
    eresepRoutes.Put("/:id", handlers.UpdateEresep)        // Update eresep by ID
    eresepRoutes.Delete("/:id", handlers.DeleteEresep)     // Delete eresep by ID
}
