package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPasswordDemo membuat hash password sederhana menggunakan SHA-256.
// Catatan penting untuk interview:
// Ini dibuat agar project demo mudah dijalankan tanpa dependency tambahan.
// Untuk production, gunakan bcrypt/argon2, bukan SHA-256 biasa.
func HashPasswordDemo(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordDemo membandingkan password input dengan hash di database.
func CheckPasswordDemo(password string, storedHash string) bool {
	return HashPasswordDemo(password) == storedHash
}
