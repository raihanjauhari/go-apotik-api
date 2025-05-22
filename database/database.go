// file: database/database.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	// Tambahkan parseTime=true dan loc=Asia%2FJakarta agar waktu DATETIME di-parse ke timezone Asia/Jakarta
	dsn := "root:@tcp(127.0.0.1:3306)/apotik?parseTime=true&loc=Asia%2FJakarta"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal konek database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database tidak bisa di-ping:", err)
	}

	log.Println("Database connected!")
}
