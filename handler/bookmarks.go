package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新增书签（验证分组归属）
func CreateBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	groupID := c.Param("id")
	groupIDInt, _ := strconv.Atoi(groupID)

	// 验证分组归属
	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", groupIDInt, user.ID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分组不存在或无权限"})
		return
	}

	var req struct {
		Name    string `json:"name" binding:"required"`
		URL     string `json:"url" binding:"required"`
		Icon    string `json:"icon"`
		BgColor string `json:"bg_color"`
		IconBg  string `json:"icon_bg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
		return
	}

	if req.Icon == "" {
		req.Icon = "🔗"
	}

	var maxSort int
	db.DB.Model(&model.Bookmark{}).Where("group_id = ?", groupIDInt).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	sort := maxSort + 1

	bookmark := model.Bookmark{
		UserID:  user.ID,
		GroupID: groupIDInt,
		Name:    req.Name,
		URL:     req.URL,
		Icon:    req.Icon,
		Sort:    sort,
		BgColor: req.BgColor,
		IconBg:  req.IconBg,
	}

	err := db.DB.Create(&bookmark).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Bookmark{
		ID:      bookmark.ID,
		Name:    req.Name,
		URL:     req.URL,
		Icon:    req.Icon,
		Sort:    sort,
		BgColor: req.BgColor,
		IconBg:  req.IconBg,
	})
}

// 更新书签（验证归属）
func UpdateBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	groupID := c.Param("id")
	bookmarkID := c.Param("bookmarkId")

	var req struct {
		Name    string `json:"name" binding:"required"`
		URL     string `json:"url" binding:"required"`
		Icon    string `json:"icon"`
		BgColor string `json:"bg_color"`
		IconBg  string `json:"icon_bg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
		return
	}

	if req.Icon == "" {
		req.Icon = "🔗"
	}

	// 直接通过 user_id + group_id 验证归属
	err := db.DB.Model(&model.Bookmark{}).
		Where("id = ? AND group_id = ? AND user_id = ?", bookmarkID, groupID, user.ID).
		Updates(map[string]interface{}{
			"name":     req.Name,
			"url":      req.URL,
			"icon":     req.Icon,
			"bg_color": req.BgColor,
			"icon_bg":  req.IconBg,
		}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	bookmarkIDInt, _ := strconv.Atoi(bookmarkID)
	c.JSON(http.StatusOK, model.Bookmark{
		ID:      bookmarkIDInt,
		Name:    req.Name,
		URL:     req.URL,
		Icon:    req.Icon,
		BgColor: req.BgColor,
		IconBg:  req.IconBg,
	})
}

// 删除书签（验证归属）
func DeleteBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	groupID := c.Param("id")
	bookmarkID := c.Param("bookmarkId")

	// 直接通过 user_id 验证归属
	err := db.DB.Where("id = ? AND group_id = ? AND user_id = ?", bookmarkID, groupID, user.ID).Delete(&model.Bookmark{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 重排序书签（验证归属）
func ReorderBookmarks(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	groupID := c.Param("id")
	groupIDInt, _ := strconv.Atoi(groupID)

	// 验证分组归属
	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", groupIDInt, user.ID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分组不存在或无权限"})
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
		err := tx.Model(&model.Bookmark{}).Where("id = ? AND user_id = ?", id, user.ID).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "排序失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
