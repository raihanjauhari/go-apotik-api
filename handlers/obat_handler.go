package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Obat struct {
	KodeObat    string  `json:"kode_obat"`
	IDUser      *string `json:"id_user,omitempty"` // bisa null
	NamaObat    string  `json:"nama_obat"`
	HargaSatuan float64 `json:"harga_satuan"`
	Stok        int     `json:"stok"`
	Deskripsi   string  `json:"deskripsi"`
}

// GET /api/obat
func GetAllObat(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT KODE_OBAT, ID_USER, NAMA_OBAT, HARGA_SATUAN, STOK, DESKRIPSI FROM OBAT")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		defer rows.Close()

		var obats []Obat
		for rows.Next() {
			var o Obat
			err := rows.Scan(&o.KodeObat, &o.IDUser, &o.NamaObat, &o.HargaSatuan, &o.Stok, &o.Deskripsi)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Error scanning data"})
			}
			obats = append(obats, o)
		}
		return c.JSON(obats)
	}
}

// GET /api/obat/:kode
func GetObatByKode(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		kode := c.Params("kode")
		var o Obat
		err := db.QueryRow("SELECT KODE_OBAT, ID_USER, NAMA_OBAT, HARGA_SATUAN, STOK, DESKRIPSI FROM OBAT WHERE KODE_OBAT = ?", kode).
			Scan(&o.KodeObat, &o.IDUser, &o.NamaObat, &o.HargaSatuan, &o.Stok, &o.Deskripsi)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).JSON(fiber.Map{"error": "Obat tidak ditemukan"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		return c.JSON(o)
	}
}

// POST /api/obat
func CreateObat(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		o := new(Obat)
		if err := c.BodyParser(o); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Insert ke DB
		_, err := db.Exec("INSERT INTO OBAT (KODE_OBAT, ID_USER, NAMA_OBAT, HARGA_SATUAN, STOK, DESKRIPSI) VALUES (?, ?, ?, ?, ?, ?)",
			o.KodeObat, o.IDUser, o.NamaObat, o.HargaSatuan, o.Stok, o.Deskripsi)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed insert data"})
		}

		return c.Status(201).JSON(o)
	}
}

// PUT /api/obat/:kode
func UpdateObat(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		kode := c.Params("kode")
		o := new(Obat)
		if err := c.BodyParser(o); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Update data
		res, err := db.Exec("UPDATE OBAT SET ID_USER=?, NAMA_OBAT=?, HARGA_SATUAN=?, STOK=?, DESKRIPSI=? WHERE KODE_OBAT=?",
			o.IDUser, o.NamaObat, o.HargaSatuan, o.Stok, o.Deskripsi, kode)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed update data"})
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Obat tidak ditemukan"})
		}

		return c.JSON(fiber.Map{"message": "Obat berhasil diupdate"})
	}
}

// DELETE /api/obat/:kode
func DeleteObat(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		kode := c.Params("kode")

		res, err := db.Exec("DELETE FROM OBAT WHERE KODE_OBAT=?", kode)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed delete data"})
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "Obat tidak ditemukan"})
		}

		return c.JSON(fiber.Map{"message": "Obat berhasil dihapus"})
	}
}
