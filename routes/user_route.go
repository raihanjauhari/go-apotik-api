package route

import (
	"go-apotik-api/handlers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	r := app.Group("/api/user")
	r.Get("/", handlers.GetAllUser)
	r.Get("/:id", handlers.GetUserByID)
	r.Post("/", handlers.CreateUser)
	r.Put("/:id", handlers.UpdateUser)
	r.Delete("/:id", handlers.DeleteUser)
	r.Post("/send-reset-code", handlers.SendResetCode)
	r.Post("/verify-reset-code", handlers.VerifyResetCode) // <-- PENTING: ini di dalam group /api/user
}
