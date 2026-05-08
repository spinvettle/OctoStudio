package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

func CustomRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {

		stackTrace := string(debug.Stack()) //获取堆栈

		slog.ErrorContext(c.Request.Context(), "system panic recovery",
			slog.Any("error", recovered),
			slog.String("stack", stackTrace),
			slog.String("path", c.Request.URL.Path),
		)
		utils.FAIL(c, http.StatusInternalServerError, "error:", "server internal error")
	})
}
