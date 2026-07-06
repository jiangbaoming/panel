package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger Gin 请求日志中间件
// 记录每个请求的方法、路径、状态码、客户端 IP 和耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		latency := time.Since(start)

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Int64("latency_ms", latency.Milliseconds()),
		}
		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		// 根据状态码选择日志级别
		level := slog.LevelInfo
		msg := "请求完成"
		if status >= 500 {
			level = slog.LevelError
			msg = "服务器错误"
		} else if status >= 400 {
			level = slog.LevelWarn
			msg = "客户端错误"
		}

		slog.LogAttrs(c.Request.Context(), level, msg, attrs...)
	}
}
