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

func GetGroups(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	var groups []model.Group
	err := db.DB.Where("user_id = ?", user.ID).Order("sort ASC").Find(&groups).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	groupIDs := make([]int, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	var bookmarks []model.Bookmark
	if len(groupIDs) > 0 {
		err = db.DB.Where("user_id = ? AND group_id IN ?", user.ID, groupIDs).Order("sort ASC").Find(&bookmarks).Error
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "查询失败", err)
			return
		}
	}

	result := make([]model.Group, 0, len(groups))
	for _, g := range groups {
		group := g
		for _, b := range bookmarks {
			if b.GroupID == g.ID {
				bookmarkCopy := b
				group.Bookmarks = append(group.Bookmarks, bookmarkCopy)
			}
		}
		result = append(result, group)
	}

	response.OK(c, result)
}

func CreateGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少名称", nil)
		return
	}

	if req.Icon == "" {
		req.Icon = "📁"
	}

	var maxSort int
	db.DB.Model(&model.Group{}).Where("user_id = ?", user.ID).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	sort := maxSort + 1

	group := model.Group{
		UserID: user.ID,
		Name:   req.Name,
		Icon:   req.Icon,
		Sort:   sort,
	}

	err := db.DB.Create(&group).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建失败", err)
		return
	}

	response.OK(c, model.Group{
		ID:        group.ID,
		UserID:    user.ID,
		Name:      req.Name,
		Icon:      req.Icon,
		Sort:      sort,
		Bookmarks: []model.Bookmark{},
	})
}

func UpdateGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少名称", nil)
		return
	}

	if req.Icon == "" {
		req.Icon = "📁"
	}

	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
		response.Error(c, http.StatusNotFound, "分组不存在或无权限", nil)
		return
	}

	err := db.DB.Model(&model.Group{}).Where("id = ? AND user_id = ?", id, user.ID).Updates(map[string]interface{}{
		"name": req.Name,
		"icon": req.Icon,
	}).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新失败", err)
		return
	}

	groupID, _ := strconv.Atoi(id)
	response.OK(c, model.Group{
		ID:   groupID,
		Name: req.Name,
		Icon: req.Icon,
	})
}

func DeleteGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	id := c.Param("id")

	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
		response.Error(c, http.StatusNotFound, "分组不存在或无权限", nil)
		return
	}

	db.DB.Where("group_id = ?", id).Delete(&model.Bookmark{})
	err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&model.Group{}).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}
	response.OK(c, gin.H{"success": true})
}

func ReorderGroups(c *gin.Context) {
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
		var group model.Group
		if err := tx.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
			tx.Rollback()
			response.Error(c, http.StatusNotFound, "分组不存在或无权限", nil)
			return
		}
		err := tx.Model(&model.Group{}).Where("id = ? AND user_id = ?", id, user.ID).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			response.Error(c, http.StatusInternalServerError, "排序失败", err)
			return
		}
	}

	tx.Commit()
	response.OK(c, gin.H{"success": true})
}
