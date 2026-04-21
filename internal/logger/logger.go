package logger

import (
	"log/slog"
	"os"
)

func InitLogger(mode string, path string) error {
	var level slog.LevelVar
	var handler slog.Handler
	if mode == "production" {
		level.Set(slog.LevelError)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: &level,
		})

	} else {
		level.Set(slog.LevelDebug)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: &level,
		})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil

}
