package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

func FAIL(c *gin.Context, httpCode int, message string, data any) {
	c.JSON(httpCode, gin.H{
		"error": message,
		"data":  data,
	})
}
