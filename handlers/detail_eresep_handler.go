package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllDetailEresep - ambil semua detail_eresep
func GetAllDetailEresep(c *fiber.Ctx) error {
    rows, err := database.DB.Query("SELECT ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, KUANTITAS, ATURAN_PAKAI, CATATAN FROM DETAIL_ERESEP")
    if err != nil {
        log.Println("Error fetching detail_eresep:", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
    }
    defer rows.Close()

    var details []models.DetailEresep
    for rows.Next() {
        var d models.DetailEresep
        var tanggalStr string
        if err := rows.Scan(&d.IDDetail, &d.IDEresep, &tanggalStr, &d.Kuantitas, &d.AturanPakai, &d.Catatan); err != nil {
            log.Println("Error scanning detail_eresep:", err)
            return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
        }
        d.TanggalEresep, _ = time.Parse("2006-01-02", tanggalStr)
        details = append(details, d)
    }

    return c.JSON(details)
}

// GetDetailEresepByID - ambil detail_eresep berdasarkan ID_DETAIL
func GetDetailEresepByID(c *fiber.Ctx) error {
    id := c.Params("id")
    var d models.DetailEresep
    var tanggalStr string

    err := database.DB.QueryRow("SELECT ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, KUANTITAS, ATURAN_PAKAI, CATATAN FROM DETAIL_ERESEP WHERE ID_DETAIL = ?", id).
        Scan(&d.IDDetail, &d.IDEresep, &tanggalStr, &d.Kuantitas, &d.AturanPakai, &d.Catatan)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Detail ERESEP not found"})
    }

    d.TanggalEresep, _ = time.Parse("2006-01-02", tanggalStr)
    return c.JSON(d)
}

// CreateDetailEresep - tambah detail_eresep baru
func CreateDetailEresep(c *fiber.Ctx) error {
    d := new(models.DetailEresep)
    if err := c.BodyParser(d); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Format tanggal ke string yyyy-mm-dd untuk MySQL
    tanggalStr := d.TanggalEresep.Format("2006-01-02")

    _, err := database.DB.Exec("INSERT INTO DETAIL_ERESEP (ID_DETAIL, ID_ERESEP, TANGGAL_ERESEP, KUANTITAS, ATURAN_PAKAI, CATATAN) VALUES (?, ?, ?, ?, ?, ?)",
        d.IDDetail, d.IDEresep, tanggalStr, d.Kuantitas, d.AturanPakai, d.Catatan)
    if err != nil {
        log.Println("Error inserting detail_eresep:", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create detail_eresep"})
    }

    return c.Status(201).JSON(d)
}

// UpdateDetailEresep - update data detail_eresep by ID_DETAIL
func UpdateDetailEresep(c *fiber.Ctx) error {
    id := c.Params("id")
    d := new(models.DetailEresep)
    if err := c.BodyParser(d); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }

    tanggalStr := d.TanggalEresep.Format("2006-01-02")

    res, err := database.DB.Exec("UPDATE DETAIL_ERESEP SET ID_ERESEP=?, TANGGAL_ERESEP=?, KUANTITAS=?, ATURAN_PAKAI=?, CATATAN=? WHERE ID_DETAIL=?",
        d.IDEresep, tanggalStr, d.Kuantitas, d.AturanPakai, d.Catatan, id)
    if err != nil {
        log.Println("Error updating detail_eresep:", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to update detail_eresep"})
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Detail ERESEP not found"})
    }

    d.IDDetail = id
    return c.JSON(d)
}

// DeleteDetailEresep - hapus detail_eresep by ID_DETAIL
func DeleteDetailEresep(c *fiber.Ctx) error {
    id := c.Params("id")

    res, err := database.DB.Exec("DELETE FROM DETAIL_ERESEP WHERE ID_DETAIL = ?", id)
    if err != nil {
        log.Println("Error deleting detail_eresep:", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to delete detail_eresep"})
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return c.Status(404).JSON(fiber.Map{"error": "Detail ERESEP not found"})
    }

    return c.SendStatus(204) // No Content
}
