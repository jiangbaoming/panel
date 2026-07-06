package handler

import (
	"net/http"
	"panel/db"
	"panel/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取用户设置
func GetSettings(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))

	var settings model.UserSettings
	err := db.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		// 不存在则创建
		newSettings := model.UserSettings{
			UserID:      userID,
			BgImage:     "",
			DisplayMode: "both",
			PageTitle:   "个人导航页",
			PageFavicon: "",
		}
		db.DB.Create(&newSettings)
		settings = newSettings
	}

	c.JSON(http.StatusOK, settings)
}

// 更新用户设置
func UpdateSettings(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))

	var req struct {
		BgImage     string `json:"bg_image"`
		DisplayMode string `json:"display_mode"`
		PageTitle   string `json:"page_title"`
		PageFavicon string `json:"page_favicon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查是否存在
	var settings model.UserSettings
	err := db.DB.Where("user_id = ?", userID).First(&settings).Error

	if err != nil {
		// 插入
		settings = model.UserSettings{
			UserID:      userID,
			BgImage:     req.BgImage,
			DisplayMode: req.DisplayMode,
			PageTitle:   req.PageTitle,
			PageFavicon: req.PageFavicon,
		}
		db.DB.Create(&settings)
	} else {
		// 更新
		updates := map[string]interface{}{}
		if req.BgImage != "" {
			updates["bg_image"] = req.BgImage
		}
		if req.DisplayMode != "" {
			updates["display_mode"] = req.DisplayMode
		}
		if req.PageTitle != "" {
			updates["page_title"] = req.PageTitle
		}
		if req.PageFavicon != "" {
			updates["page_favicon"] = req.PageFavicon
		}
		if len(updates) > 0 {
			db.DB.Model(&model.UserSettings{}).Where("user_id = ?", userID).Updates(updates)
		}
	}

	// 返回最新设置
	db.DB.Where("user_id = ?", userID).First(&settings)
	c.JSON(http.StatusOK, settings)
}
