package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/middlewares"
)

func Router() *gin.Engine {
	var r *gin.Engine
	if gin.Mode() == "debug" {
		r = gin.Default()
	} else {
		r = gin.New()
	}
	r.Use(middlewares.TraceID())
	r.Use(middlewares.Auth())
	{
		v1Group := r.Group("/v1")
		{
			v1Group.GET("models", func(c *gin.Context) {
				a := []int{}
				fmt.Print(a[0])
				c.JSON(http.StatusOK, gin.H{
					"msg": "ok",
				})

			})
			v1Group.POST("/chat/completions")
			v1Group.POST("/completions")
			v1Group.POST("/response")
			v1Group.POST("/embeddings")
		}
	}
	return r
}
