package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/proxy"
)

func Router(relay *proxy.RelayHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/backend-api/codex/responses", relay.Relay)

	return r
}
