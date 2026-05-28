package utils

import "github.com/gin-gonic/gin"

// SuccessResponse membuat format response sukses yang konsisten.
func SuccessResponse(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse membuat format response error yang konsisten.
func ErrorResponse(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"error":   err,
	})
}
