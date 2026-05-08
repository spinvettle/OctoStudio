package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spinvettle/OctoStudio/internal/consts"
	"gorm.io/gorm/logger"
)

type TraceHandler struct {
	slog.Handler
}

// Handle 拦截日志记录动作
func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	// 尝试从 context 中拿出 trace_id
	if traceID, ok := ctx.Value(consts.CtxKeyTraceID).(string); ok && traceID != "" {
		// 如果有，就自动加上一个 attribute
		r.AddAttrs(slog.String(string(consts.CtxKeyTraceID), traceID))
	}
	return h.Handler.Handle(ctx, r)
}

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
	traceHandler := &TraceHandler{handler}
	logger := slog.New(traceHandler)
	slog.SetDefault(logger)

	return nil

}

type MyGormLogger struct {
	LogLevel logger.LogLevel
}

func (l *MyGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &MyGormLogger{LogLevel: level} // 返回新实例
}

func (l *MyGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		slog.InfoContext(ctx, "GORM Info", "message", fmt.Sprintf(msg, args...))
	}
}

func (l *MyGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		slog.WarnContext(ctx, "GORM Warn", "message", fmt.Sprintf(msg, args...))
	}
}

func (l *MyGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		slog.ErrorContext(ctx, "GORM Error", "message", fmt.Sprintf(msg, args...))
	}
}

func (l *MyGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 处理错误打 Error 日志
	if err != nil && l.LogLevel >= logger.Error {
		slog.ErrorContext(ctx, "GORM SQL Error",
			"sql", sql,
			"RowsAffected", rows,
			"elapsed", elapsed,
			"error", err,
		)
		return
	}

	// 处理慢查询：假设超过 200ms
	if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		slog.WarnContext(ctx, "GORM SLOW SQL",
			"sql", sql,
			"RowsAffected", rows,
			"elapsed", elapsed,
		)
		return
	}

	// 处理正常记录：打 Info 日志
	if l.LogLevel >= logger.Info {
		slog.InfoContext(ctx, "GORM SQL Trace",
			"sql", sql,
			"RowsAffected", rows,
			"elapsed", elapsed,
		)
	}
}
