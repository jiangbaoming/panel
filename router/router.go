package router

import (
	"os"
	"panel/config"
	"panel/handler"
	"panel/middleware"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup 注册所有路由并返回 gin 引擎
func Setup() *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 100 << 20 // 100MB
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// 静态文件：上传目录
	r.Static("/uploads", config.UploadDir)

	// 公开路由（无需认证）
	r.POST("/api/login", handler.Login)

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())
	{
		setupUserRoutes(api)
		setupGroupRoutes(api)
		setupBookmarkRoutes(api)
		setupPinnedRoutes(api)
		setupImageRoutes(api)
		setupSettingRoutes(api)
		setupSearchRoutes(api)
	}

	// 生产环境托管前端静态文件（SPA fallback）
	setupFrontendStatic(r)

	return r
}

// 用户管理路由
func setupUserRoutes(api *gin.RouterGroup) {
	api.GET("/users", middleware.AdminRequired(), handler.GetUsers)
	api.POST("/users", middleware.AdminRequired(), handler.CreateUser)
	api.DELETE("/users/:id", middleware.AdminRequired(), handler.DeleteUser)
	api.PATCH("/users/me", handler.UpdateMe)
	api.POST("/change-password", handler.ChangePassword)
}

// 分组管理路由
func setupGroupRoutes(api *gin.RouterGroup) {
	api.GET("/groups", handler.GetGroups)
	api.POST("/groups", handler.CreateGroup)
	api.PUT("/groups/:id", handler.UpdateGroup)
	api.DELETE("/groups/:id", handler.DeleteGroup)
	api.PATCH("/groups/reorder", handler.ReorderGroups)
}

// 书签管理路由
func setupBookmarkRoutes(api *gin.RouterGroup) {
	api.POST("/groups/:id/bookmarks", handler.CreateBookmark)
	api.PUT("/groups/:id/bookmarks/:bookmarkId", handler.UpdateBookmark)
	api.DELETE("/groups/:id/bookmarks/:bookmarkId", handler.DeleteBookmark)
	api.PATCH("/groups/:id/bookmarks/reorder", handler.ReorderBookmarks)
}

// 常驻书签路由
func setupPinnedRoutes(api *gin.RouterGroup) {
	api.GET("/pinned", handler.GetPinned)
	api.PATCH("/pinned/:id/pin", handler.TogglePin)
	api.PATCH("/pinned/reorder", handler.ReorderPinned)
}

// 图片管理路由
func setupImageRoutes(api *gin.RouterGroup) {
	api.POST("/images/upload", handler.UploadImage)
	api.POST("/images/upload-zip", handler.UploadZip)
	api.GET("/images", handler.GetImages)
	api.DELETE("/images/:id", middleware.AdminRequired(), handler.DeleteImage)
	api.DELETE("/images/clear/:category", middleware.AdminRequired(), handler.ClearCategory)
}

// 用户设置路由
func setupSettingRoutes(api *gin.RouterGroup) {
	api.GET("/settings/:userId", handler.GetSettings)
	api.PUT("/settings/:userId", handler.UpdateSettings)
}

// 搜索路由
func setupSearchRoutes(api *gin.RouterGroup) {
	api.GET("/search", handler.SearchBookmarks)
}

// 前端静态文件托管
func setupFrontendStatic(r *gin.Engine) {
	distDir := getDistDir()
	if distDir == "" {
		return
	}

	r.StaticFile("/favicon.svg", filepath.Join(distDir, "favicon.svg"))
	r.Static("/assets", filepath.Join(distDir, "assets"))
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(distDir, "index.html"))
	})
}

// 获取前端构建产物目录
func getDistDir() string {
	// 优先使用环境变量
	distDir := os.Getenv("DIST_DIR")
	if distDir != "" {
		return distDir
	}

	// 尝试多个可能的位置
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	candidates := []string{
		filepath.Join(execDir, "dist"),
		filepath.Join(execDir, "..", "dist"),
		filepath.Join(execDir, "..", "..", "dist"),
		"./dist",
		"../dist",
	}

	for _, dir := range candidates {
		if _, err := os.Stat(filepath.Join(dir, "index.html")); err == nil {
			return dir
		}
	}

	return ""
}
