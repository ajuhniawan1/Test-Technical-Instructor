package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const idleSessionPrefix = "session_idle:"

type idleSessionResponseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *idleSessionResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *idleSessionResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}

	return w.body.Write(data)
}

func (w *idleSessionResponseWriter) WriteString(s string) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = 200
	}

	return w.body.WriteString(s)
}

// RegisterIdleSessionAfterLogin dipasang di route login.
// Setelah login sukses, middleware ini membaca token dari response JSON,
// lalu menyimpan token hash ke Redis dengan TTL idle timeout.
func RegisterIdleSessionAfterLogin(redisClient *redis.Client, idleTTL time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		// Simpan writer asli dari Gin.
		originalWriter := c.Writer

		// Pakai writer custom untuk menangkap response login.
		bufferWriter := &idleSessionResponseWriter{
			ResponseWriter: originalWriter,
			body:           bytes.NewBuffer(nil),
			statusCode:     200,
		}

		c.Writer = bufferWriter

		// Jalankan handler login.
		c.Next()

		statusCode := bufferWriter.statusCode
		responseBody := bufferWriter.body.Bytes()

		// Kalau login sukses, ambil token dari response body lalu simpan idle session ke Redis.
		if statusCode >= 200 && statusCode < 300 {
			token := extractTokenFromLoginResponse(responseBody)

			if token != "" {
				ttl := calculateIdleSessionTTL(token, idleTTL)

				if ttl > 0 {
					err := CreateIdleSession(c.Request.Context(), redisClient, token, ttl)
					if err != nil {
						log.Println("failed to create idle session:", err)
					}
				}
			}
		}

		// Kembalikan writer asli supaya response benar-benar dikirim ke client/Postman.
		c.Writer = originalWriter
		c.Writer.WriteHeader(statusCode)

		if len(responseBody) > 0 {
			_, _ = c.Writer.Write(responseBody)
		}
	}
}

// IdleSessionMiddleware dipasang setelah AuthMiddleware.
// Middleware ini mengecek apakah session masih aktif di Redis.
// Jika aktif, TTL diperpanjang lagi.
func IdleSessionMiddleware(redisClient *redis.Client, idleTTL time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "Redis tidak tersedia untuk validasi session",
			})
			return
		}

		token := GetBearerToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Authorization token tidak ditemukan",
			})
			return
		}

		ctx := c.Request.Context()
		key := buildIdleSessionKey(token)

		exists, err := redisClient.Exists(ctx, key).Result()
		if err != nil {
			log.Println("idle session check error:", err)

			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "Gagal validasi session",
			})
			return
		}

		if exists == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Session expired karena tidak ada aktivitas. Silakan login ulang.",
			})
			return
		}

		ttl := calculateIdleSessionTTL(token, idleTTL)
		if ttl <= 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token sudah expired. Silakan login ulang.",
			})
			return
		}

		if err := redisClient.Expire(ctx, key, ttl).Err(); err != nil {
			log.Println("idle session refresh error:", err)

			c.AbortWithStatusJSON(503, gin.H{
				"success": false,
				"message": "Gagal memperpanjang session",
			})
			return
		}

		c.Next()
	}
}

func CreateIdleSession(ctx context.Context, redisClient *redis.Client, token string, ttl time.Duration) error {
	key := buildIdleSessionKey(token)
	return redisClient.Set(ctx, key, "active", ttl).Err()
}

func DeleteIdleSession(ctx context.Context, redisClient *redis.Client, token string) error {
	key := buildIdleSessionKey(token)
	return redisClient.Del(ctx, key).Err()
}

func buildIdleSessionKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return idleSessionPrefix + hex.EncodeToString(hash[:])
}

func extractTokenFromLoginResponse(responseBody []byte) string {
	var payload map[string]interface{}

	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return ""
	}

	if token, ok := payload["token"].(string); ok {
		return token
	}

	if accessToken, ok := payload["access_token"].(string); ok {
		return accessToken
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return ""
	}

	if token, ok := data["token"].(string); ok {
		return token
	}

	if accessToken, ok := data["access_token"].(string); ok {
		return accessToken
	}

	return ""
}

func calculateIdleSessionTTL(token string, idleTTL time.Duration) time.Duration {
	tokenRemainingTTL := getJWTRemainingDurationForIdleSession(token)

	if tokenRemainingTTL <= 0 {
		return 0
	}

	if tokenRemainingTTL < idleTTL {
		return tokenRemainingTTL
	}

	return idleTTL
}

func getJWTRemainingDurationForIdleSession(tokenString string) time.Duration {
	claims := jwt.MapClaims{}

	parser := jwt.NewParser()

	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		// Kalau token tidak bisa dibaca exp-nya, fallback ke 30 menit.
		return 30 * time.Minute
	}

	expValue, ok := claims["exp"]
	if !ok {
		return 30 * time.Minute
	}

	var expUnix int64

	switch value := expValue.(type) {
	case float64:
		expUnix = int64(value)
	case int64:
		expUnix = value
	default:
		return 30 * time.Minute
	}

	return time.Until(time.Unix(expUnix, 0))
}