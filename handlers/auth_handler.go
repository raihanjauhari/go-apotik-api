package handlers

import (
	"fmt"
	"go-apotik-api/database"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SendResetCode - kirim kode reset password ke email user
func SendResetCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	// Cek email ada di database
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user WHERE email = ?)", req.Email).Scan(&exists)
	if err != nil {
		log.Println("Error cek email:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Server error"})
	}
	if !exists {
		return c.Status(404).JSON(fiber.Map{"error": "Email tidak ditemukan"})
	}

	// Generate kode 6 digit
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Kirim email
	err = sendEmail(req.Email, code)
	if err != nil {
		log.Println("Gagal kirim email:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengirim email verifikasi."})
	}

	// Hapus kode sebelumnya
	_, err = database.DB.Exec(
		"DELETE FROM reset_codes WHERE id_user = (SELECT id_user FROM user WHERE email = ?)",
		req.Email,
	)
	if err != nil {
		log.Println("Gagal hapus kode sebelumnya:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengatur ulang kode sebelumnya."})
	}

	// Simpan kode baru
	_, err = database.DB.Exec(
		"INSERT INTO reset_codes (id_user, code) VALUES ((SELECT id_user FROM user WHERE email = ?), ?)",
		req.Email, code,
	)
	if err != nil {
		log.Println("Gagal simpan kode:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan kode"})
	}

	log.Println("SendResetCode sukses untuk email:", req.Email)

	return c.JSON(fiber.Map{
		"message": "Kode verifikasi telah dikirim ke email Anda",
	})
}

// Fungsi kirim email pakai SMTP Gmail
func sendEmail(to string, code string) error {
	from := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	if from == "" || pass == "" {
		return fmt.Errorf("Environment variable EMAIL_USER atau EMAIL_PASS belum diset")
	}

	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: Kode Verifikasi Reset Password\r\n" +
			"\r\n" +
			"Kode verifikasi Anda adalah: " + code + "\r\n",
	)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Error saat kirim email ke %s: %v\n", to, err)
		return err
	}

	log.Println("Email berhasil dikirim ke", to)
	return nil
}
