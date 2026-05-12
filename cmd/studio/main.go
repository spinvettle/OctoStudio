package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/db"
	"github.com/spinvettle/OctoStudio/internal/gateway/channel"
	"github.com/spinvettle/OctoStudio/internal/gateway/relay"
	"github.com/spinvettle/OctoStudio/internal/httpclient"
	"github.com/spinvettle/OctoStudio/internal/logger"
	"github.com/spinvettle/OctoStudio/internal/router"
	"gorm.io/gorm"
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
	var err error
	var DB *gorm.DB
	Init()
	port := config.Port
	// mode := config.GlobalConfig.ServerConfig.Mode
	// gin.SetMode(conf)
	if DB, err = db.InitDB(config.DSN); err != nil {
		panic(err)
	}
	client := httpclient.NewRelayClient(time.Minute * 10)
	channelSvc := channel.NewChannelService(DB, client)
	handler := relay.NewRelayHandler(channelSvc)
	server := router.Router(handler)

	slog.Info("Run server", "port", port, "mode", config.Mode)
	err = server.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("Server Run Error", "error", err)
	}
}
