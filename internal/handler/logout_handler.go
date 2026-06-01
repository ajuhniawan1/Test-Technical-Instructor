package handler

import (
	"log"

	"assignment-platform/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Logout memasukkan JWT yang sedang dipakai ke Redis blacklist.
// Jadi token yang sama tidak bisa dipakai lagi setelah logout.
func Logout(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := middleware.GetBearerToken(c)

		if token == "" {
			c.JSON(401, gin.H{
				"success": false,
				"message": "Authorization token tidak ditemukan",
			})
			return
		}

		err := middleware.AddTokenToBlacklist(c.Request.Context(), redisClient, token)
		if err != nil {
			log.Println("failed to blacklist token:", err)

			c.JSON(500, gin.H{
				"success": false,
				"message": "Gagal logout",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Logout berhasil",
		})
	}
}