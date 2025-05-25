package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetAllDilayani - ambil semua data dilayani
func GetAllDilayani(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT ID_PENDAFTARAN, ID_DOKTER FROM DILAYANI")
	if err != nil {
		log.Println("Error fetching dilayani:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var list []models.Dilayani
	for rows.Next() {
		var d models.Dilayani
		if err := rows.Scan(&d.IDPendaftaran, &d.IDDokter); err != nil {
			log.Println("Error scanning dilayani:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		list = append(list, d)
	}

	return c.JSON(list)
}

// GetDilayaniByID - ambil data dilayani berdasarkan ID_PENDAFTARAN
func GetDilayaniByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var d models.Dilayani

	err := database.DB.QueryRow("SELECT ID_PENDAFTARAN, ID_DOKTER FROM DILAYANI WHERE ID_PENDAFTARAN = ?", id).
		Scan(&d.IDPendaftaran, &d.IDDokter)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	return c.JSON(d)
}

// CreateDilayani - tambah data dilayani baru
func CreateDilayani(c *fiber.Ctx) error {
	d := new(models.Dilayani)
	if err := c.BodyParser(d); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	_, err := database.DB.Exec("INSERT INTO DILAYANI (ID_PENDAFTARAN, ID_DOKTER) VALUES (?, ?)",
		d.IDPendaftaran, d.IDDokter)
	if err != nil {
		log.Println("Error inserting dilayani:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create data"})
	}

	return c.Status(201).JSON(d)
}

// UpdateDilayani - update data dilayani by ID_PENDAFTARAN
func UpdateDilayani(c *fiber.Ctx) error {
	id := c.Params("id")
	d := new(models.Dilayani)
	if err := c.BodyParser(d); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	res, err := database.DB.Exec("UPDATE DILAYANI SET ID_DOKTER=? WHERE ID_PENDAFTARAN=?",
		d.IDDokter, id)
	if err != nil {
		log.Println("Error updating dilayani:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update data"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	d.IDPendaftaran = id
	return c.JSON(d)
}

// DeleteDilayani - hapus data dilayani by ID_PENDAFTARAN
func DeleteDilayani(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM DILAYANI WHERE ID_PENDAFTARAN = ?", id)
	if err != nil {
		log.Println("Error deleting dilayani:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete data"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Data not found"})
	}

	return c.SendStatus(204) // No Content
}
