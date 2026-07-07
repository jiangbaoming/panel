package middleware

import (
	"bytes"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter 包装 gin.ResponseWriter，捕获响应体内容
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger Gin 请求日志中间件
// 记录每个请求的方法、路径、状态码、客户端 IP 和耗时
// 4xx/5xx 响应时额外记录响应体内容，方便排查错误
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 包装响应写入器，捕获响应体
		blw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = blw

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

		// 4xx/5xx 时追加响应体内容
		if status >= 400 {
			body := blw.body.String()
			if body != "" {
				attrs = append(attrs, slog.String("body", body))
			}
		}

		// 追加 gin 框架记录的 error（如有）
		for _, e := range c.Errors {
			attrs = append(attrs, slog.String("gin_error", e.Err.Error()))
		}

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
