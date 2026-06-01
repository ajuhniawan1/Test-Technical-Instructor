package handler

import (
	"log"

	"assignment-platform/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Logout memasukkan JWT ke Redis blacklist dan menghapus idle session.
// Jadi token tidak bisa dipakai lagi setelah logout.
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

		if err := middleware.DeleteIdleSession(c.Request.Context(), redisClient, token); err != nil {
			log.Println("failed to delete idle session:", err)
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Logout berhasil",
		})
	}
}