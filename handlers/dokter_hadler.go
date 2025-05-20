package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetAllDokter - ambil semua dokter
func GetAllDokter(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT ID_DOKTER, NAMA_DOKTER, POLI, FOTO_DOKTER FROM DOKTER")
	if err != nil {
		log.Println("Error fetching dokter:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var dokters []models.Dokter
	for rows.Next() {
		var d models.Dokter
		if err := rows.Scan(&d.ID, &d.Nama, &d.Poli, &d.FotoDokter); err != nil {
			log.Println("Error scanning dokter:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		dokters = append(dokters, d)
	}

	return c.JSON(dokters)
}

// GetDokterByID - ambil dokter berdasarkan ID
func GetDokterByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var d models.Dokter

	err := database.DB.QueryRow("SELECT ID_DOKTER, NAMA_DOKTER, POLI, FOTO_DOKTER FROM DOKTER WHERE ID_DOKTER = ?", id).
		Scan(&d.ID, &d.Nama, &d.Poli, &d.FotoDokter)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Dokter not found"})
	}

	return c.JSON(d)
}

// CreateDokter - tambah dokter baru
func CreateDokter(c *fiber.Ctx) error {
	d := new(models.Dokter)
	if err := c.BodyParser(d); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	_, err := database.DB.Exec("INSERT INTO DOKTER (ID_DOKTER, NAMA_DOKTER, POLI, FOTO_DOKTER) VALUES (?, ?, ?, ?)",
		d.ID, d.Nama, d.Poli, d.FotoDokter)
	if err != nil {
		log.Println("Error inserting dokter:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create dokter"})
	}

	return c.Status(201).JSON(d)
}

// UpdateDokter - update data dokter by ID
func UpdateDokter(c *fiber.Ctx) error {
	id := c.Params("id")
	d := new(models.Dokter)
	if err := c.BodyParser(d); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	res, err := database.DB.Exec("UPDATE DOKTER SET NAMA_DOKTER=?, POLI=?, FOTO_DOKTER=? WHERE ID_DOKTER=?",
		d.Nama, d.Poli, d.FotoDokter, id)
	if err != nil {
		log.Println("Error updating dokter:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update dokter"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Dokter not found"})
	}

	d.ID = id
	return c.JSON(d)
}

// DeleteDokter - hapus dokter by ID
func DeleteDokter(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM DOKTER WHERE ID_DOKTER = ?", id)
	if err != nil {
		log.Println("Error deleting dokter:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete dokter"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Dokter not found"})
	}

	return c.SendStatus(204) // No Content
}
