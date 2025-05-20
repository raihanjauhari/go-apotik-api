package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllEresep - ambil semua data resep
func GetAllEresep(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT ID_ERESEP, ID_PENDAFTARAN, STATUS FROM ERESEP")
	if err != nil {
		log.Println("Error fetching eresep:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var ereseps []models.Eresep
	for rows.Next() {
		var e models.Eresep
		if err := rows.Scan(&e.IDEresep, &e.IDPendaftaran, &e.Status); err != nil {
			log.Println("Error scanning eresep:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		ereseps = append(ereseps, e)
	}

	return c.JSON(ereseps)
}

// GetEresepByID - ambil data resep berdasarkan ID_ERESEP
func GetEresepByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.Eresep

	err := database.DB.QueryRow("SELECT ID_ERESEP, ID_PENDAFTARAN, STATUS FROM ERESEP WHERE ID_ERESEP = ?", id).
		Scan(&e.IDEresep, &e.IDPendaftaran, &e.Status)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Eresep not found"})
	}

	return c.JSON(e)
}

// CreateEresep - tambah data resep baru
func CreateEresep(c *fiber.Ctx) error {
	e := new(models.Eresep)
	if err := c.BodyParser(e); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Generate ID unik
	e.IDEresep = uuid.New().String()

	// Validasi status
	validStatus := map[string]bool{
		"Menunggu Pembayaran": true,
		"Sudah Bayar":         true,
		"Diproses":            true,
		"Selesai":             true,
	}
	if !validStatus[e.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status value"})
	}

	_, err := database.DB.Exec("INSERT INTO ERESEP (ID_ERESEP, ID_PENDAFTARAN, STATUS) VALUES (?, ?, ?)",
		e.IDEresep, e.IDPendaftaran, e.Status)
	if err != nil {
		log.Println("Error inserting eresep:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create eresep"})
	}

	return c.Status(201).JSON(e)
}

// UpdateEresep - update data resep by ID_ERESEP
func UpdateEresep(c *fiber.Ctx) error {
	id := c.Params("id")
	e := new(models.Eresep)
	if err := c.BodyParser(e); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi status jika diupdate
	validStatus := map[string]bool{
		"Menunggu Pembayaran": true,
		"Sudah Bayar":         true,
		"Diproses":            true,
		"Selesai":             true,
	}
	if e.Status != "" && !validStatus[e.Status] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid status value"})
	}

	// Update data, hanya jika field diisi
	res, err := database.DB.Exec(
		"UPDATE ERESEP SET ID_PENDAFTARAN = IF(? = '', ID_PENDAFTARAN, ?), STATUS = IF(? = '', STATUS, ?) WHERE ID_ERESEP = ?",
		e.IDPendaftaran, e.IDPendaftaran, e.Status, e.Status, id)
	if err != nil {
		log.Println("Error updating eresep:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update eresep"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Eresep not found"})
	}

	e.IDEresep = id
	return c.JSON(e)
}

// DeleteEresep - hapus data resep by ID_ERESEP
func DeleteEresep(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM ERESEP WHERE ID_ERESEP = ?", id)
	if err != nil {
		log.Println("Error deleting eresep:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete eresep"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Eresep not found"})
	}

	return c.SendStatus(204) // No Content
}
