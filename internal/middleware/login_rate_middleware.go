package middleware

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var loginRateLimiterScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])

if current == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end

return current
`)

// LoginRateLimiter membatasi percobaan login berdasarkan IP.
// Contoh: max 5 request per 1 menit.
func LoginRateLimiter(redisClient *redis.Client, maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		ctx := context.Background()

		ip := c.ClientIP()
		key := "rate_limit:login:" + ip

		count, err := loginRateLimiterScript.Run(
			ctx,
			redisClient,
			[]string{key},
			int(window.Seconds()),
		).Int64()

		if err != nil {
			log.Println("rate limiter redis error:", err)

			// Kalau Redis error, request tetap boleh jalan
			// agar aplikasi tidak mati hanya karena Redis bermasalah.
			c.Next()
			return
		}

		remaining := maxRequests - count
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(maxRequests, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if count > maxRequests {
			ttl, _ := redisClient.TTL(ctx, key).Result()

			c.AbortWithStatusJSON(429, gin.H{
				"success": false,
				"message": "Terlalu banyak percobaan login. Silakan coba lagi nanti.",
				"retry_after_seconds": int(ttl.Seconds()),
			})
			return
		}

		c.Next()
	}
}