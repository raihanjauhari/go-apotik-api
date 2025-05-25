package handlers

import (
	"fmt"
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// formatTanggalIndo format time.Time ke "17 Mei 2025"
func formatTanggalIndo(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	bulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return fmt.Sprintf("%02d %s %d", t.Day(), bulan[t.Month()-1], t.Year())
}

// DetailEresepResponse untuk response API dengan tanggal string
type DetailEresepResponse struct {
	IDDetail      string `json:"id_detail"`
	IDEresep      string `json:"id_eresep"`
	TanggalEresep string `json:"tanggal_eresep"`
	Catatan       string `json:"catatan"`
}

// GetAllDetailEresep ambil semua data detail_eresep
func GetAllDetailEresep(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, CATATAN FROM DETAIL_ERESEP")
	if err != nil {
		log.Println("Error fetching detail_eresep:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var details []DetailEresepResponse

	for rows.Next() {
		var d models.DetailEresep
		// Scan langsung tanggal ke time.Time (pastikan DB driver mendukung)
		if err := rows.Scan(&d.IDDetail, &d.IDEresep, &d.TanggalEresep, &d.Catatan); err != nil {
			log.Println("Error scanning detail_eresep:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read data"})
		}

		details = append(details, DetailEresepResponse{
			IDDetail:      d.IDDetail,
			IDEresep:      d.IDEresep,
			TanggalEresep: formatTanggalIndo(d.TanggalEresep),
			Catatan:       d.Catatan,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to iterate data"})
	}

	return c.JSON(details)
}

// GetDetailEresepByID ambil data detail_eresep berdasarkan ID_DETAIL
func GetDetailEresepByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var d models.DetailEresep

	err := database.DB.QueryRow("SELECT ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, CATATAN FROM DETAIL_ERESEP WHERE ID_DETAIL = ?", id).
		Scan(&d.IDDetail, &d.IDEresep, &d.TanggalEresep, &d.Catatan)
	if err != nil {
		log.Println("Error fetching detail_eresep by id:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Detail ERESEP not found"})
	}

	detailResp := DetailEresepResponse{
		IDDetail:      d.IDDetail,
		IDEresep:      d.IDEresep,
		TanggalEresep: formatTanggalIndo(d.TanggalEresep),
		Catatan:       d.Catatan,
	}

	return c.JSON(detailResp)
}

// CreateDetailEresep tambah data baru detail_eresep
func CreateDetailEresep(c *fiber.Ctx) error {
	type inputCreate struct {
		IDDetail      string `json:"id_detail"`
		IDEresep      string `json:"id_eresep"`
		TanggalEresep string `json:"tanggal_eresep"` // format "YYYY-MM-DD"
		Catatan       string `json:"catatan"`
	}

	var input inputCreate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	tanggal, err := time.Parse("2006-01-02", input.TanggalEresep)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tanggal_eresep format, gunakan YYYY-MM-DD"})
	}

	_, err = database.DB.Exec("INSERT INTO DETAIL_ERESEP (ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, CATATAN) VALUES (?, ?, ?, ?)",
		input.IDDetail, input.IDEresep, tanggal, input.Catatan)
	if err != nil {
		log.Println("Error inserting detail_eresep:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create detail_eresep"})
	}

	detailResp := DetailEresepResponse{
		IDDetail:      input.IDDetail,
		IDEresep:      input.IDEresep,
		TanggalEresep: formatTanggalIndo(tanggal),
		Catatan:       input.Catatan,
	}

	return c.Status(fiber.StatusCreated).JSON(detailResp)
}

// UpdateDetailEresep update data detail_eresep berdasarkan ID_DETAIL
func UpdateDetailEresep(c *fiber.Ctx) error {
	id := c.Params("id")

	type inputUpdate struct {
		IDEresep      string `json:"id_eresep"`
		TanggalEresep string `json:"tanggal_eresep"`
		Catatan       string `json:"catatan"`
	}

	var input inputUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	tanggal, err := time.Parse("2006-01-02", input.TanggalEresep)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tanggal_eresep format, gunakan YYYY-MM-DD"})
	}

	res, err := database.DB.Exec("UPDATE DETAIL_ERESEP SET ID_ERESEP=?, TANGGAL_ERESEP=?, CATATAN=? WHERE ID_DETAIL=?",
		input.IDEresep, tanggal, input.Catatan, id)
	if err != nil {
		log.Println("Error updating detail_eresep:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update detail_eresep"})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update detail_eresep"})
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Detail ERESEP not found"})
	}

	detailResp := DetailEresepResponse{
		IDDetail:      id,
		IDEresep:      input.IDEresep,
		TanggalEresep: formatTanggalIndo(tanggal),
		Catatan:       input.Catatan,
	}

	return c.JSON(detailResp)
}

// DeleteDetailEresep hapus data detail_eresep berdasarkan ID_DETAIL
func DeleteDetailEresep(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM DETAIL_ERESEP WHERE ID_DETAIL = ?", id)
	if err != nil {
		log.Println("Error deleting detail_eresep:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete detail_eresep"})
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete detail_eresep"})
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Detail ERESEP not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
