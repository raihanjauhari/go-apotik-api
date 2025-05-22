// file: handlers/auth_handler.go
package handlers

import (
	"fmt"
	"go-apotik-api/database"

	"github.com/gofiber/fiber/v2"
)

func VerifyResetCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	fmt.Println("Verifying reset code for email:", req.Email, "code:", req.Code)

	var idUser string
	err := database.DB.QueryRow("SELECT id_user FROM user WHERE email = ?", req.Email).Scan(&idUser)
	if err != nil {
		fmt.Println("Email tidak ditemukan:", req.Email)
		return c.Status(404).JSON(fiber.Map{"error": "Email tidak ditemukan"})
	}

	var savedCode string
	err = database.DB.QueryRow("SELECT code FROM reset_codes WHERE id_user = ? AND code = ?", idUser, req.Code).Scan(&savedCode)
	if err != nil {
		fmt.Println("Kode tidak ditemukan atau belum diminta untuk user:", idUser, "code:", req.Code)
		return c.Status(404).JSON(fiber.Map{"error": "Kode tidak ditemukan atau belum diminta"})
	}

	// Tidak lagi memeriksa expired time
	return c.JSON(fiber.Map{"message": "Kode verifikasi benar"})
}
