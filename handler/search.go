package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"panel/response"

	"github.com/gin-gonic/gin"
)

func SearchBookmarks(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未认证", nil)
		return
	}

	q := c.Query("q")
	if q == "" {
		response.OK(c, []model.SearchResult{})
		return
	}

	var results []model.SearchResult
	err := db.DB.Table("bookmarks b").
		Select("b.id, b.group_id, b.name, b.url, b.icon, g.name as groupName").
		Joins("JOIN groups_t g ON b.group_id = g.id").
		Where("b.user_id = ? AND (b.name LIKE ? OR b.url LIKE ?)", user.ID, "%"+q+"%", "%"+q+"%").
		Scan(&results).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	response.OK(c, results)
}
