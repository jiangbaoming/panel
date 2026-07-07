package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"panel/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPinned(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
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
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	if pinned == nil {
		pinned = []model.PinnedBookmark{}
	}
	response.OK(c, pinned)
}

func TogglePin(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	id := c.Param("id")

	var bookmark model.Bookmark
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&bookmark).Error; err != nil {
		response.Error(c, http.StatusNotFound, "书签不存在或无权限", nil)
		return
	}

	var req struct {
		Pinned bool `json:"pinned"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少参数", nil)
		return
	}

	pinnedVal := 0
	if req.Pinned {
		pinnedVal = 1
	}

	err := db.DB.Model(&model.Bookmark{}).Where("id = ?", id).Update("pinned", pinnedVal).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	response.OK(c, gin.H{
		"id":     id,
		"pinned": pinnedVal,
	})
}

func ReorderPinned(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	var req model.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		response.Error(c, http.StatusBadRequest, "缺少 ids 数组", nil)
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
		var bookmark model.Bookmark
		if err := tx.Where("id = ? AND user_id = ?", id, user.ID).First(&bookmark).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusNotFound, "书签不存在或无权限", nil)
			return
		}
		err := tx.Model(&model.Bookmark{}).Where("id = ?", id).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "排序失败", err)
			return
		}
	}

	tx.Commit()
	response.OK(c, gin.H{"success": true})
}
