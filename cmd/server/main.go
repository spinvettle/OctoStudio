package main

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/config"
	"github.com/spinvettle/OctoStudio/internal/logger"
	"github.com/spinvettle/OctoStudio/internal/router"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	// 初始化日志
	err := logger.InitLogger(
		config.GlobalConfig.LoggingConfig.Level,
		config.GlobalConfig.LoggingConfig.Output,
		config.GlobalConfig.LoggingConfig.FileDir,
	)
	if err != nil {
		panic(err)
	}
	// TODO llama-server 资源检测
	slog.Error("缺少bin/llama-server服务")
	port := config.GlobalConfig.ServerConfig.Port
	mode := config.GlobalConfig.ServerConfig.Mode

	gin.SetMode(mode)
	server := router.Router()

	slog.Info("Run server", "port", port, "mdoe", gin.Mode)
	err = server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Server Run Error", "error", err)
	}
}
