package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisLockKeyFunc func(c *gin.Context) string

var releaseLockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`)

// RedisLock mencegah request yang sama diproses bersamaan.
// Cocok untuk mencegah race condition saat submit assignment.
func RedisLock(redisClient *redis.Client, ttl time.Duration, keyFunc RedisLockKeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "Redis tidak tersedia untuk locking",
			})
			return
		}
	// untuk generate key yang unik berdasarkan request, misal berdasarkan user id dan assignment id.
		lockKey := keyFunc(c)
// Jika key tidak valid, bisa jadi karena parameter yang dibutuhkan tidak ada. Abort dengan 400.
		if lockKey == "" {
			c.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"message": "Lock key tidak valid",
			})
			return
		}
// Generate token unik untuk lock ini, agar saat release kita bisa pastikan hanya yang punya token ini yang bisa release locknya.
		ctx := c.Request.Context()
		lockToken := generateLockToken()

		// SET key value NX EX
		// NX = hanya set kalau key belum ada
		// TTL = lock otomatis hilang kalau ada error/panic
		ok, err := redisClient.SetNX(ctx, lockKey, lockToken, ttl).Result()
		if err != nil {
			log.Println("redis lock error:", err)

			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "Gagal membuat lock Redis",
			})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(409, gin.H{
				"success": false,
				"message": "Request yang sama sedang diproses. Silakan tunggu sebentar.",
			})
			return
		}

		// Lepas lock setelah handler selesai.
		// defer releaseRedisLock(ctx, redisClient, lockKey, lockToken)
		defer func() {
			log.Println("RELEASE REDIS LOCK:", lockKey)

			err := releaseLockScript.Run(ctx, redisClient, []string{lockKey}, lockToken).Err()
			if err != nil {
				log.Println("REDIS RELEASE LOCK ERROR:", err)
			}
		}()
		
		c.Next()
	}
}

// SubmissionCreateLockKey membuat lock berdasarkan assignment_id dan user_id.
// Jadi 1 talent tidak bisa submit assignment yang sama secara bersamaan.
func SubmissionCreateLockKey(c *gin.Context) string {
	assignmentID := c.Param("assignmentId")
	userID := getUserIDFromContext(c)

	if assignmentID == "" || userID == "" {
		return ""
	}

	return fmt.Sprintf("lock:submission:create:assignment:%s:user:%s", assignmentID, userID)
}

// SubmissionUpdateLockKey untuk resubmit berdasarkan submission id dan user id.
func SubmissionUpdateLockKey(c *gin.Context) string {
	submissionID := c.Param("id")
	userID := getUserIDFromContext(c)

	if submissionID == "" || userID == "" {
		return ""
	}

	return fmt.Sprintf("lock:submission:update:submission:%s:user:%s", submissionID, userID)
}

func getUserIDFromContext(c *gin.Context) string {
	possibleKeys := []string{
		"user_id",
		"userID",
		"UserID",
		"id",
	}

	for _, key := range possibleKeys {
		value, exists := c.Get(key)
		if exists {
			return fmt.Sprint(value)
		}
	}

	return ""
}

func generateLockToken() string {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(bytes)
}

func releaseRedisLock(ctx context.Context, redisClient *redis.Client, lockKey string, lockToken string) {
	err := releaseLockScript.Run(ctx, redisClient, []string{lockKey}, lockToken).Err()
	if err != nil {
		log.Println("failed to release redis lock:", err)
	}
}