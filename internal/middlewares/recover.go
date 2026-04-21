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

		slog.Error("System Panic Recovery",
			"error", recovered,
			"stack", stackTrace,
			"path", c.Request.URL.Path,
		)
		utils.FAIL(c, http.StatusInternalServerError, "error:内部错误", nil)
	})
}
