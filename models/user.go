package models

type User struct {
	IDUser   string `json:"id_user" db:"id_user"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"` // Simpan hashed password
	NamaUser string `json:"nama_user" db:"nama_user"`
	Role     string `json:"role" db:"role"`           // "admin" atau "petugas"
	FotoUser string `json:"foto_user" db:"foto_user"` // nama file foto user
}
