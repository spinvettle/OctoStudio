package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/proxy"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/backend-api/codex/responses", proxy.Relay)

	return r
}
