package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetAllMemunculkan - ambil semua data memunculkan
func GetAllMemunculkan(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT KODE_OBAT, ID_ERESEP, ID_DETAIL, KUANTITAS, ATURAN_PAKAI FROM memunculkan")
	if err != nil {
		log.Println("Error fetching memunculkan:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var memunculkans []models.Memunculkan
	for rows.Next() {
		var m models.Memunculkan
		if err := rows.Scan(&m.KodeObat, &m.IDEresep, &m.IDDetail, &m.Kuantitas, &m.AturanPakai); err != nil {
			log.Println("Error scanning memunculkan:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		memunculkans = append(memunculkans, m)
	}

	return c.JSON(memunculkans)
}

// GetMemunculkanByIDs - ambil data berdasarkan KODE_OBAT dan ID_ERESEP
func GetMemunculkanByIDs(c *fiber.Ctx) error {
	kodeObat := c.Params("kode_obat")
	idEresep := c.Params("id_eresep")

	var m models.Memunculkan
	err := database.DB.QueryRow("SELECT KODE_OBAT, ID_ERESEP, ID_DETAIL, KUANTITAS, ATURAN_PAKAI FROM memunculkan WHERE KODE_OBAT = ? AND ID_ERESEP = ?", kodeObat, idEresep).
		Scan(&m.KodeObat, &m.IDEresep, &m.IDDetail, &m.Kuantitas, &m.AturanPakai)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	return c.JSON(m)
}

// CreateMemunculkan - tambah data baru
func CreateMemunculkan(c *fiber.Ctx) error {
	m := new(models.Memunculkan)
	if err := c.BodyParser(m); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	_, err := database.DB.Exec("INSERT INTO memunculkan (KODE_OBAT, ID_ERESEP, ID_DETAIL, KUANTITAS, ATURAN_PAKAI) VALUES (?, ?, ?, ?, ?)",
		m.KodeObat, m.IDEresep, m.IDDetail, m.Kuantitas, m.AturanPakai)
	if err != nil {
		log.Println("Error inserting memunculkan:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create data"})
	}

	return c.Status(201).JSON(m)
}

// UpdateMemunculkan - update data berdasarkan KODE_OBAT dan ID_ERESEP (misal update ID_DETAIL, KUANTITAS, ATURAN_PAKAI)
func UpdateMemunculkan(c *fiber.Ctx) error {
	kodeObat := c.Params("kode_obat")
	idEresep := c.Params("id_eresep")

	m := new(models.Memunculkan)
	if err := c.BodyParser(m); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	res, err := database.DB.Exec("UPDATE memunculkan SET ID_DETAIL=?, KUANTITAS=?, ATURAN_PAKAI=? WHERE KODE_OBAT = ? AND ID_ERESEP = ?",
		m.IDDetail, m.Kuantitas, m.AturanPakai, kodeObat, idEresep)
	if err != nil {
		log.Println("Error updating memunculkan:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update data"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	m.KodeObat = kodeObat
	m.IDEresep = idEresep
	return c.JSON(m)
}

// DeleteMemunculkan - hapus data berdasarkan KODE_OBAT dan ID_ERESEP
func DeleteMemunculkan(c *fiber.Ctx) error {
	kodeObat := c.Params("kode_obat")
	idEresep := c.Params("id_eresep")

	res, err := database.DB.Exec("DELETE FROM memunculkan WHERE KODE_OBAT = ? AND ID_ERESEP = ?", kodeObat, idEresep)
	if err != nil {
		log.Println("Error deleting memunculkan:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete data"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	return c.SendStatus(204)
}
