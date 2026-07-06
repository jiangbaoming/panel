package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取当前用户的分组（包括书签）
func GetGroups(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 查询当前用户的分组
	var groups []model.Group
	err := db.DB.Where("user_id = ?", user.ID).Order("sort ASC").Find(&groups).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 收集分组ID
	groupIDs := make([]int, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	// 查询这些分组下的书签
	var bookmarks []model.Bookmark
	if len(groupIDs) > 0 {
		err = db.DB.Where("user_id = ? AND group_id IN ?", user.ID, groupIDs).Order("sort ASC").Find(&bookmarks).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
	}

	// 组装数据
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

	c.JSON(http.StatusOK, result)
}

// 新增分组（归属当前用户）
func CreateGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少名称"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, model.Group{
		ID:        group.ID,
		UserID:    user.ID,
		Name:      req.Name,
		Icon:      req.Icon,
		Sort:      sort,
		Bookmarks: []model.Bookmark{},
	})
}

// 更新分组（仅允许操作自己的分组）
func UpdateGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	id := c.Param("id")
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少名称"})
		return
	}

	if req.Icon == "" {
		req.Icon = "📁"
	}

	// 验证分组归属
	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分组不存在或无权限"})
		return
	}

	err := db.DB.Model(&model.Group{}).Where("id = ? AND user_id = ?", id, user.ID).Updates(map[string]interface{}{
		"name": req.Name,
		"icon": req.Icon,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	groupID, _ := strconv.Atoi(id)
	c.JSON(http.StatusOK, model.Group{
		ID:   groupID,
		Name: req.Name,
		Icon: req.Icon,
	})
}

// 删除分组（仅允许操作自己的分组）
func DeleteGroup(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	id := c.Param("id")

	// 验证分组归属
	var group model.Group
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分组不存在或无权限"})
		return
	}

	// 先删除关联书签
	db.DB.Where("group_id = ?", id).Delete(&model.Bookmark{})
	// 再删除分组
	err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&model.Group{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 重排序分组（仅允许操作自己的分组）
func ReorderGroups(c *gin.Context) {
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
		// 验证归属
		var group model.Group
		if err := tx.Where("id = ? AND user_id = ?", id, user.ID).First(&group).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "分组不存在或无权限"})
			return
		}
		err := tx.Model(&model.Group{}).Where("id = ? AND user_id = ?", id, user.ID).Update("sort", i).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "排序失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
