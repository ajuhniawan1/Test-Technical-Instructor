package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type cacheResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *cacheResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *cacheResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ResponseCache dipakai untuk cache response GET.
// Contoh: GET /api/v1/classes
func ResponseCache(redisClient *redis.Client, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		// Cache hanya untuk GET
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		ctx := context.Background()

		// Dibuat aman agar cache user A tidak kebaca user B.
		scope := getCacheScope(c)

		cacheKey := fmt.Sprintf(
			"cache:%s:%s:%s",
			scope,
			c.Request.Method,
			c.Request.URL.RequestURI(),
		)

		// Cek Redis dulu
		cached, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			c.Header("X-Cache", "HIT")
			c.Data(200, "application/json; charset=utf-8", []byte(cached))
			c.Abort()
			return
		}

		// Kalau belum ada cache, lanjut ke handler asli
		writer := &cacheResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}

		c.Writer = writer
		c.Header("X-Cache", "MISS")

		c.Next()

		statusCode := c.Writer.Status()
		contentType := c.Writer.Header().Get("Content-Type")

		// Simpan ke Redis hanya kalau response sukses dan JSON
		if statusCode == 200 && strings.Contains(contentType, "application/json") {
			err := redisClient.Set(ctx, cacheKey, writer.body.String(), ttl).Err()
			if err != nil {
				log.Println("failed to set cache:", err)
			}
		}
	}
}

// InvalidateCache dipakai setelah POST/PUT/DELETE supaya cache lama dihapus.
func InvalidateCache(redisClient *redis.Client, patterns ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if redisClient == nil {
			return
		}

		statusCode := c.Writer.Status()

		// Hapus cache hanya kalau create/update/delete sukses
		if statusCode < 200 || statusCode >= 300 {
			return
		}

		ctx := context.Background()

		for _, pattern := range patterns {
			err := deleteRedisKeysByPattern(ctx, redisClient, pattern)
			if err != nil {
				log.Println("failed to invalidate cache:", err)
			}
		}
	}
}

func deleteRedisKeysByPattern(ctx context.Context, redisClient *redis.Client, pattern string) error {
	iter := redisClient.Scan(ctx, 0, pattern, 100).Iterator()

	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())

		if len(keys) >= 100 {
			if err := redisClient.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			keys = keys[:0]
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return redisClient.Del(ctx, keys...).Err()
	}

	return nil
}

func getCacheScope(c *gin.Context) string {
	// Kalau auth middleware kamu menyimpan user_id di context,
	// cache akan berdasarkan user id.
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	if userID, exists := c.Get("userID"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	// Fallback aman: pakai hash Authorization token.
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		hash := sha256.Sum256([]byte(authHeader))
		return "auth:" + hex.EncodeToString(hash[:])[:16]
	}

	return "public"
}