package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"assignment-platform/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectMySQL membuka koneksi ke database MySQL.
// Kita memakai database/sql agar query lebih eksplisit dan mudah dijelaskan saat interview.
func ConnectMySQL(cfg config.Config) *sql.DB {
	// parseTime=true penting agar field DATETIME MySQL bisa discan ke time.Time di Go.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to open database connection: ", err)
	}

	// Connection pool sederhana.
	// Ini membantu mengontrol jumlah koneksi aplikasi ke MySQL.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping dipakai untuk memastikan koneksi benar-benar bisa digunakan.
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	log.Println("connected to MySQL successfully")
	return db
}
