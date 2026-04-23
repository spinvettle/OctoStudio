package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/spinvettle/OctoStudio/internal/consts"
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
