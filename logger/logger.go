package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

// Level 日志级别
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var levelMap = map[string]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
}

// Init 初始化全局日志记录器
// level: debug / info / warn / error，空值默认为 info
// json: true 输出 JSON 格式，false 输出文本格式（默认）
func Init(level string, json bool) {
	lvl := slog.LevelInfo
	if v, ok := levelMap[level]; ok {
		lvl = v
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: lvl,
		// 记录调用源文件和行号
		AddSource: lvl <= slog.LevelDebug,
	}

	if json {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// NewTextHandler 创建一个彩色文本 Handler
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &textHandler{
		Handler: slog.NewTextHandler(w, opts),
		levels: levelColor{
			slog.LevelDebug: "36", // 青色
			slog.LevelInfo:  "32", // 绿色
			slog.LevelWarn:  "33", // 黄色
			slog.LevelError: "31", // 红色
		},
	}
}

type levelColor map[slog.Level]string

type textHandler struct {
	slog.Handler
	levels levelColor
}

func (h *textHandler) Handle(ctx context.Context, r slog.Record) error {
	// 从 Attrs 中提取 message，避免重复输出 level/msg
	return h.Handler.Handle(ctx, r)
}

// Fatal 记录错误日志并终止程序（替代 log.Fatalf）
func Fatal(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	slog.Error(msg)
	os.Exit(1)
}

// Debug 便捷方法
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Info 便捷方法
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn 便捷方法
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error 便捷方法
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Elapsed 返回耗时日志 Attr
func Elapsed(d time.Duration) slog.Attr {
	return slog.String("elapsed", fmt.Sprintf("%dms", d.Milliseconds()))
}
