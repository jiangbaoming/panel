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

func CreateBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	groupID := c.Param("id")
	groupIDInt, _ := strconv.Atoi(groupID)

	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", groupIDInt, user.ID).First(&group).Error; err != nil {
		response.Error(c, http.StatusNotFound, "分组不存在或无权限", nil)
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
		response.Error(c, http.StatusBadRequest, "缺少必要参数", nil)
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
		response.Error(c, http.StatusInternalServerError, "创建失败", err)
		return
	}

	response.OK(c, model.Bookmark{
		ID:      bookmark.ID,
		Name:    req.Name,
		URL:     req.URL,
		Icon:    req.Icon,
		Sort:    sort,
		BgColor: req.BgColor,
		IconBg:  req.IconBg,
	})
}

func UpdateBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
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
		response.Error(c, http.StatusBadRequest, "缺少必要参数", nil)
		return
	}

	if req.Icon == "" {
		req.Icon = "🔗"
	}

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
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	bookmarkIDInt, _ := strconv.Atoi(bookmarkID)
	response.OK(c, model.Bookmark{
		ID:      bookmarkIDInt,
		Name:    req.Name,
		URL:     req.URL,
		Icon:    req.Icon,
		BgColor: req.BgColor,
		IconBg:  req.IconBg,
	})
}

func DeleteBookmark(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	groupID := c.Param("id")
	bookmarkID := c.Param("bookmarkId")

	err := db.DB.Where("id = ? AND group_id = ? AND user_id = ?", bookmarkID, groupID, user.ID).Delete(&model.Bookmark{}).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}

	response.OK(c, gin.H{"success": true})
}

func ReorderBookmarks(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	groupID := c.Param("id")
	groupIDInt, _ := strconv.Atoi(groupID)

	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", groupIDInt, user.ID).First(&group).Error; err != nil {
		response.Error(c, http.StatusNotFound, "分组不存在或无权限", nil)
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
		err := tx.Model(&model.Bookmark{}).Where("id = ? AND user_id = ?", id, user.ID).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "排序失败", err)
			return
		}
	}

	tx.Commit()
	response.OK(c, gin.H{"success": true})
}
