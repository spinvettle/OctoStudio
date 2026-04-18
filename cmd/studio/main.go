package main

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	// "github.com/spinvettle/OctoStudio/config"
	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/router"
)

func main() {
	if err := config.LoadConfig("./config.yaml"); err != nil {
		panic(err)
	}

	port := config.GlobalConfig.Port
	// mode := config.GlobalConfig.ServerConfig.Mode

	server := router.Router()

	slog.Info("Run server", "port", port, "mdoe", gin.Mode)
	err := server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Server Run Error", "error", err)
	}
}
