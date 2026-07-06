package main

import (
	"log"
	"panel/config"
	"panel/db"
	"panel/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	db.Init()

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 注册路由
	r := router.Setup()

	// 启动服务
	log.Printf("后端服务已启动: http://localhost:%s", config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
