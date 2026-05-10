package db

import (
	"errors"
	"strings"

	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
	var DB *gorm.DB
	if dsn == "" {
		return nil, errors.New("db dsn is incorrect")
	}
	var err error
	var level gormLogger.LogLevel
	lowerLevel := strings.ToLower(config.LogLevel)
	switch lowerLevel {
	case "error":
		level = gormLogger.Error
	case "info":
		level = gormLogger.Info
	case "warn":
		level = gormLogger.Warn
	default:
		level = gormLogger.Silent
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &logger.MyGormLogger{LogLevel: level},
	})
	if err != nil {
		return nil, err
	}
	return DB, nil
}
