package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Result 统一响应结构
type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Error 返回错误响应，自动将原始错误注入 gin 上下文供日志中间件捕获
func Error(c *gin.Context, status int, msg string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.JSON(status, Result{
		Code:    status,
		Message: msg,
		Data:    nil,
	})
}

// OK 返回成功响应
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Result{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}
