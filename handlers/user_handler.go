package handlers

import (
	"go-apotik-api/database"
	"go-apotik-api/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// helper untuk hashing password
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// GetAllUser - ambil semua user
func GetAllUser(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id_user, email, password, nama_user, role, foto_user FROM user")
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

	err := database.DB.QueryRow("SELECT id_user, email, password, nama_user, role, foto_user FROM user WHERE id_user = ?", id).
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

	// Validasi password minimal 6 karakter
	if len(u.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
	}

	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	u.Password = hashedPassword

	u.IDUser = uuid.New().String()

	_, err = database.DB.Exec("INSERT INTO user (id_user, email, password, nama_user, role, foto_user) VALUES (?, ?, ?, ?, ?, ?)",
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

	// Jika password ada dan tidak kosong, hash dulu
	if u.Password != "" {
		if len(u.Password) < 6 {
			return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
		}
		hashedPassword, err := hashPassword(u.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		u.Password = hashedPassword
	}

	res, err := database.DB.Exec("UPDATE user SET email=?, password=?, nama_user=?, role=?, foto_user=? WHERE id_user=?",
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

	res, err := database.DB.Exec("DELETE FROM user WHERE id_user = ?", id)
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

// UpdatePassword - update password user by ID
func UpdatePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	var input struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if len(input.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters"})
	}

	hashed, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	res, err := database.DB.Exec("UPDATE user SET password=? WHERE id_user=?", hashed, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.SendStatus(200)
}

func LoginUser(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    var user models.User
    err := database.DB.QueryRow("SELECT id_user, email, password, nama_user, role, foto_user FROM user WHERE email = ?", input.Email).
        Scan(&user.IDUser, &user.Email, &user.Password, &user.NamaUser, &user.Role, &user.FotoUser)

    if err != nil {
        // User tidak ditemukan
        log.Println("Login error - user not found:", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
    }

    // Cocokkan password dari input dengan hash di database
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        // Password salah
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
    }

    // Password cocok, hapus password sebelum kirim response
    user.Password = ""

    return c.JSON(user)
}

