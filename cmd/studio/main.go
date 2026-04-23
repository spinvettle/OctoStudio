package main

import (
	"fmt"
	"log/slog"

	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/logger"
	"github.com/spinvettle/OctoStudio/internal/proxy/codexProxy"
	"github.com/spinvettle/OctoStudio/internal/router"
)

func Init() {
	if err := config.LoadConfig("./config.yaml"); err != nil {
		panic(err)
	}
	if err := logger.InitLogger(config.Mode, config.LogFile); err != nil {
		panic(err)
	}
	codexProxy.InitCodexProxy()

}

func main() {

	Init()
	port := config.Port
	// mode := config.GlobalConfig.ServerConfig.Mode

	server := router.Router()

	slog.Info("Run server", "port", port, "mode", config.Mode)
	err := server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Server Run Error", "error", err)
	}
}
