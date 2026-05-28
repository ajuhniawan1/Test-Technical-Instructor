package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menyimpan konfigurasi aplikasi.
// Dipisahkan ke struct agar konfigurasi tidak tersebar di banyak file.
type Config struct {
	AppPort   string
	JWTSecret string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
}

// Load membaca konfigurasi dari .env.
// Jika .env tidak ada, aplikasi tetap mencoba membaca dari environment variable.
func Load() Config {
	// godotenv.Load() tidak dibuat fatal agar aplikasi tetap bisa jalan di server
	// yang biasanya memakai environment variable langsung.
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "change_this_secret_key"),
		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASSWORD", ""),
		DBName:    getEnv("DB_NAME", "bootcamp_assignment"),
	}
}

// getEnv mengambil nilai dari environment variable.
// Jika kosong, maka memakai default value.
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
