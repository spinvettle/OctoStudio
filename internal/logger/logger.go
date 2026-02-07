package logger

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func InitLogger(levelStr, output, fileDir string) error {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level, AddSource: true}

	var writer io.Writer
	writers := []io.Writer{}

	if output == "stdout" || output == "both" {
		writers = append(writers, os.Stdout)
	}

	if output == "file" || output == "both" {
		logFile := filepath.Join(fileDir, time.Now().Format("2006-01-02")+".log")
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return errors.New("Failed to create log directory:" + fileDir)

		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return errors.New("Failed to open log file" + err.Error())
		} else {
			writers = append(writers, file)
		}
	}

	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = os.Stdout
	}

	var handler slog.Handler
	if level == slog.LevelDebug {
		handler = slog.NewTextHandler(writer, opts)
	} else {
		handler = slog.NewJSONHandler(writer, opts)
	}

	slog.SetDefault(slog.New(handler))
	return nil
}
