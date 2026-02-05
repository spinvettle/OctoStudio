package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/middlewares"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Auth())
	{
		v1Group := r.Group("/v1")
		{
			v1Group.GET("models")
			v1Group.POST("/chat/completions")
			v1Group.POST("/completions")
			v1Group.POST("/response")
			v1Group.POST("/embeddings")
		}
	}
	return r
}
