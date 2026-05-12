package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/gateway/relay"
	"github.com/spinvettle/OctoStudio/internal/middlewares"
)

func Router(handler *relay.RelayHandler) *gin.Engine {
	if config.Mode == "release" || config.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middlewares.TraceID())
	r.Use(middlewares.AccessLog())
	r.Use(middlewares.CustomRecovery())

	r.POST("/backend-api/codex/responses", handler.Relay)
	// v1Group := r.Group("/v1")
	// {
	// 	v1Group.POST("/chat/completions", func(ctx *gin.Context) {
	// 		num, _ := strconv.Atoi(ctx.Query("a"))
	// 		fmt.Println(1 / num)
	// 	})
	// }

	return r
}
