package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const jwtBlacklistPrefix = "jwt_blacklist:"

// TokenBlacklistMiddleware mengecek apakah JWT sudah pernah logout.
// Kalau token ada di Redis blacklist, request ditolak.
func TokenBlacklistMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			c.Next()
			return
		}

		token := GetBearerToken(c)

		// Kalau token kosong, biarkan AuthMiddleware yang menolak.
		if token == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		key := buildJWTBlacklistKey(token)

		exists, err := redisClient.Exists(ctx, key).Result()
		if err != nil {
			log.Println("redis blacklist check error:", err)

			// Fail-open: kalau Redis error, request tetap lanjut.
			// Kalau ingin lebih ketat, bisa return 503 di sini.
			c.Next()
			return
		}

		if exists > 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token sudah logout. Silakan login ulang.",
			})
			return
		}

		c.Next()
	}
}

// AddTokenToBlacklist menyimpan JWT ke Redis sampai waktu expired JWT habis.
func AddTokenToBlacklist(ctx context.Context, redisClient *redis.Client, token string) error {
	if redisClient == nil {
		return errors.New("redis client is nil")
	}

	if token == "" {
		return errors.New("token is empty")
	}

	ttl := getTokenRemainingTTL(token)

	// Kalau token sudah expired, tidak perlu disimpan.
	if ttl <= 0 {
		return nil
	}

	key := buildJWTBlacklistKey(token)

	return redisClient.Set(ctx, key, "1", ttl).Err()
}

// GetBearerToken mengambil token dari header Authorization.
func GetBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

// buildJWTBlacklistKey membuat key Redis dari hash token.
// Token asli tidak disimpan langsung supaya lebih aman.
func buildJWTBlacklistKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return jwtBlacklistPrefix + hex.EncodeToString(hash[:])
}

// getTokenRemainingTTL membaca sisa waktu hidup JWT dari claim exp.
func getTokenRemainingTTL(tokenString string) time.Duration {
	claims := jwt.MapClaims{}

	parser := jwt.NewParser()

	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		log.Println("failed to parse jwt exp:", err)

		// Fallback sesuai JWT kamu yang biasanya 24 jam.
		return 24 * time.Hour
	}

	expValue, ok := claims["exp"]
	if !ok {
		return 24 * time.Hour
	}

	var expUnix int64

	switch value := expValue.(type) {
	case float64:
		expUnix = int64(value)
	case int64:
		expUnix = value
	default:
		return 24 * time.Hour
	}

	expTime := time.Unix(expUnix, 0)

	return time.Until(expTime)
}