package middleware

import (
	"net/http"
	"strings"

	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memastikan request memiliki JWT token yang valid.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Header harus berbentuk: Authorization: Bearer <token>
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "missing bearer token")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err.Error())
			c.Abort()
			return
		}

		// Simpan data user ke context Gin.
		// Handler/service bisa mengambilnya untuk mengetahui user yang login.
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
