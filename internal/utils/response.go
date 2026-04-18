package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, GinResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func FAIL(c *gin.Context, httpCode int, message string, data any) {
	c.JSON(httpCode, GinResponse{
		Code:    0,
		Message: message,
		Data:    data,
	})
}
