package main

import (
	"os"
	"panel/config"
	"panel/db"
	"panel/logger"
	"panel/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志（环境变量 LOG_LEVEL=debug/info/warn/error, LOG_JSON=true 启用 JSON 格式）
	logLevel := os.Getenv("LOG_LEVEL")
	logJSON := os.Getenv("LOG_JSON") == "true"
	logger.Init(logLevel, logJSON)

	// 初始化配置
	config.Init()

	// 初始化数据库
	db.Init()

	// 设置 Gin 模式
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// 注册路由
	r := router.Setup()

	// 启动服务
	logger.Info("后端服务已启动", "port", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		logger.Fatal("启动失败", "error", err)
	}
}
