package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetAllUser - ambil semua user
func GetAllUser(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id_user, email, password, nama_user, role, foto_user FROM USER")
	if err != nil {
		log.Println("Error fetching user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.IDUser, &u.Email, &u.Password, &u.NamaUser, &u.Role, &u.FotoUser); err != nil {
			log.Println("Error scanning user:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read data"})
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

// GetUserByID - ambil user berdasarkan ID
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var u models.User

	err := database.DB.QueryRow("SELECT id_user, email, password, nama_user, role, foto_user FROM USER WHERE id_user = ?", id).
		Scan(&u.IDUser, &u.Email, &u.Password, &u.NamaUser, &u.Role, &u.FotoUser)
	if err != nil {
		log.Println("User not found:", err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(u)
}

// CreateUser - tambah user baru
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	u.IDUser = uuid.New().String()

	// TODO: hash password sebelum simpan
	// u.Password = hashPassword(u.Password)

	_, err := database.DB.Exec("INSERT INTO USER (id_user, email, password, nama_user, role, foto_user) VALUES (?, ?, ?, ?, ?, ?)",
		u.IDUser, u.Email, u.Password, u.NamaUser, u.Role, u.FotoUser)
	if err != nil {
		log.Println("Error inserting user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(201).JSON(u)
}

// UpdateUser - update data user by ID
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	u := new(models.User)
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	res, err := database.DB.Exec("UPDATE USER SET email=?, password=?, nama_user=?, role=?, foto_user=? WHERE id_user=?",
		u.Email, u.Password, u.NamaUser, u.Role, u.FotoUser, id)
	if err != nil {
		log.Println("Error updating user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	u.IDUser = id
	return c.JSON(u)
}

// DeleteUser - hapus user by ID
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := database.DB.Exec("DELETE FROM USER WHERE id_user = ?", id)
	if err != nil {
		log.Println("Error deleting user:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.SendStatus(204) // No Content
}
