package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"

	"github.com/gin-gonic/gin"
)

// 搜索当前用户的书签
func SearchBookmarks(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []model.SearchResult{})
		return
	}

	var results []model.SearchResult
	err := db.DB.Table("bookmarks b").
		Select("b.id, b.group_id, b.name, b.url, b.icon, g.name as groupName").
		Joins("JOIN groups_t g ON b.group_id = g.id").
		Where("b.user_id = ? AND (b.name LIKE ? OR b.url LIKE ?)", user.ID, "%"+q+"%", "%"+q+"%").
		Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, results)
}
