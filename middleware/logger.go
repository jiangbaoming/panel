package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger Gin 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 放行
		c.Next()

		// 收集信息
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		latency := time.Since(start)
		err := c.Errors.ByType(gin.ErrorTypeAny).String()

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
		if err != "" {
			attrs = append(attrs, slog.String("error", err))
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}
		slog.LogAttrs(c.Request.Context(), level, "", attrs...)
	}
}
