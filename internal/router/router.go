package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/proxy/codexProxy"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.POST("/backend-api/codex/responses", codexProxy.CodexRelay)
	v1Group := r.Group("/v1")
	{
		v1Group.POST("/chat/completions")
	}

	return r
}
