package middleware

import (
	"net/http"

	"assignment-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi endpoint agar hanya role tertentu yang boleh mengakses.
// Contoh: RequireRole("admin", "trainer") berarti hanya admin/trainer yang boleh masuk.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "role not found in token")
			c.Abort()
			return
		}

		role := roleValue.(string)
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "your role is not allowed to access this endpoint")
		c.Abort()
	}
}
