package handlers

import (
	"go-apotik-api/database" // ganti utils jadi database supaya konsisten seperti dokter
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetAllPasien(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT ID_PENDAFTARAN, UMUR, DIAGNOSA, BERAT_BADAN, FOTO_PASIEN FROM PASIEN")
	if err != nil {
		log.Println("Error fetching pasien:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pasien"})
	}
	defer rows.Close()

	var pasienList []models.Pasien
	for rows.Next() {
		var p models.Pasien
		if err := rows.Scan(&p.IDPendaftaran, &p.Umur, &p.Diagnosa, &p.BeratBadan, &p.FotoPasien); err != nil {
			log.Println("Error scanning pasien:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read pasien data"})
		}
		pasienList = append(pasienList, p)
	}

	return c.JSON(pasienList)
}

func GetPasienByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Pasien

	err := database.DB.QueryRow("SELECT ID_PENDAFTARAN, UMUR, DIAGNOSA, BERAT_BADAN, FOTO_PASIEN FROM PASIEN WHERE ID_PENDAFTARAN = ?", id).
		Scan(&p.IDPendaftaran, &p.Umur, &p.Diagnosa, &p.BeratBadan, &p.FotoPasien)
	if err != nil {
		log.Println("Pasien tidak ditemukan:", err)
		return c.Status(404).JSON(fiber.Map{"error": "Pasien not found"})
	}

	return c.JSON(p)
}

func CreatePasien(c *fiber.Ctx) error {
	p := new(models.Pasien)

	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Cek apakah sudah ada pasien dengan ID sama
	var exists string
	err := database.DB.QueryRow("SELECT ID_PENDAFTARAN FROM PASIEN WHERE ID_PENDAFTARAN = ?", p.IDPendaftaran).Scan(&exists)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID_PENDAFTARAN sudah ada"})
	}

	_, err = database.DB.Exec("INSERT INTO PASIEN (ID_PENDAFTARAN, UMUR, DIAGNOSA, BERAT_BADAN, FOTO_PASIEN) VALUES (?, ?, ?, ?, ?)",
		p.IDPendaftaran, p.Umur, p.Diagnosa, p.BeratBadan, p.FotoPasien)
	if err != nil {
		log.Println("Error inserting pasien:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create pasien"})
	}

	return c.Status(201).JSON(p)
}

func UpdatePasien(c *fiber.Ctx) error {
	id := c.Params("id")
	p := new(models.Pasien)

	if err := c.BodyParser(p); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Pastikan id path dan body sama
	if id != p.IDPendaftaran {
		return c.Status(400).JSON(fiber.Map{"error": "ID in URL and body must be the same"})
	}

	res, err := database.DB.Exec("UPDATE PASIEN SET UMUR=?, DIAGNOSA=?, BERAT_BADAN=?, FOTO_PASIEN=? WHERE ID_PENDAFTARAN=?",
		p.Umur, p.Diagnosa, p.BeratBadan, p.FotoPasien, id)
	if err != nil {
		log.Println("Error updating pasien:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update pasien"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pasien not found"})
	}

	p.IDPendaftaran = id
	return c.JSON(p)
}

func DeletePasien(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM PASIEN WHERE ID_PENDAFTARAN = ?", id)
	if err != nil {
		log.Println("Error deleting pasien:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete pasien"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Pasien not found"})
	}

	return c.SendStatus(204) // No Content, seperti dokter handler
}
