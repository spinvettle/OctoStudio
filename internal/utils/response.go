package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/consts"
)

type GinResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"traceID,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, GinResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: c.GetString(consts.CtxKeyTraceID),
	})
}

func ERROR(c *gin.Context, code uint, data any) {
	c.JSON(http.StatusOK, GinResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: c.GetString(consts.CtxKeyTraceID),
	})
}
func BadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, GinResponse{
		Code:    400,
		Message: err,
		TraceID: c.GetString(consts.CtxKeyTraceID),
	})
}
