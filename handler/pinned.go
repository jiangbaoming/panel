package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取当前用户的常驻书签
func GetPinned(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var pinned []model.PinnedBookmark
	err := db.DB.Table("bookmarks b").
		Select("b.id, b.name, b.url, b.icon, b.bg_color, b.icon_bg, g.name as group_name").
		Joins("JOIN groups g ON b.group_id = g.id").
		Where("b.pinned = ? AND b.user_id = ?", 1, user.ID).
		Order("b.sort ASC").
		Scan(&pinned).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	if pinned == nil {
		pinned = []model.PinnedBookmark{}
	}
	c.JSON(http.StatusOK, pinned)
}

// 切换书签常驻状态（验证归属）
func TogglePin(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	id := c.Param("id")

	// 直接通过 user_id 验证归属
	var bookmark model.Bookmark
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&bookmark).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "书签不存在或无权限"})
		return
	}

	var req struct {
		Pinned bool `json:"pinned"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}

	pinnedVal := 0
	if req.Pinned {
		pinnedVal = 1
	}

	err := db.DB.Model(&model.Bookmark{}).Where("id = ?", id).Update("pinned", pinnedVal).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"pinned": pinnedVal,
	})
}

// 重排序常驻书签（验证归属）
func ReorderPinned(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req model.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 ids 数组"})
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, idStr := range req.IDs {
		id, _ := strconv.Atoi(idStr)
		// 直接通过 user_id 验证归属
		var bookmark model.Bookmark
		if err := tx.Where("id = ? AND user_id = ?", id, user.ID).First(&bookmark).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "书签不存在或无权限"})
			return
		}
		err := tx.Model(&model.Bookmark{}).Where("id = ?", id).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "排序失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
