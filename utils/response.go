package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(500).JSON(fiber.Map{
		"error": err.Error(),
	})
}
