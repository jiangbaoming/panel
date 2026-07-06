package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// Level 日志级别常量（对应环境变量 LOG_LEVEL）
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
// level: debug / info / warn / error，空值默认 info
// json: true 输出 JSON 格式，false 输出文本格式（默认）
func Init(level string, json bool) {
	lvl := slog.LevelInfo
	if v, ok := levelMap[level]; ok {
		lvl = v
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: lvl <= slog.LevelDebug, // debug 级别记录调用源文件和行号
	}

	var handler slog.Handler
	if json {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// Fatal 记录错误日志并终止程序
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

// Debug 输出调试日志
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Info 输出信息日志
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// Warn 输出警告日志
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// Error 输出错误日志
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Elapsed 返回耗时日志属性
func Elapsed(d fmt.Stringer) slog.Attr {
	return slog.String("elapsed", d.String())
}
