package main

import (
	"fmt"
	"log/slog"

	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/db"
	"github.com/spinvettle/OctoStudio/internal/logger"
	"github.com/spinvettle/OctoStudio/internal/router"
)

func Init() {
	if err := config.LoadConfig("./config.yaml"); err != nil {
		panic(err)
	}
	if err := logger.InitLogger(config.Mode, config.LogFile); err != nil {
		panic(err)
	}
	if _, err := db.InitDB(config.DSN); err != nil {
		panic(err)
	}

}

func main() {

	Init()
	port := config.Port
	// mode := config.GlobalConfig.ServerConfig.Mode
	// gin.SetMode(conf)
	server := router.Router()

	slog.Info("Run server", "port", port, "mode", config.Mode)
	err := server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Server Run Error", "error", err)
	}
}
