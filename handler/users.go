package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"
	"panel/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	var users []model.User
	err := db.DB.Select("id, username, avatar, role").Find(&users).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	response.OK(c, users)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Avatar   string `json:"avatar"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少用户名或密码", nil)
		return
	}

	if req.Avatar == "" {
		req.Avatar = "👤"
	}
	if req.Role == "" {
		req.Role = "guest"
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Username: req.Username,
		Password: string(hash),
		Avatar:   req.Avatar,
		Role:     req.Role,
	}

	err := db.DB.Create(&user).Error
	if err != nil {
		response.Error(c, http.StatusBadRequest, "用户名已存在", nil)
		return
	}

	response.OK(c, gin.H{
		"id":       user.ID,
		"username": req.Username,
		"avatar":   req.Avatar,
		"role":     req.Role,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := db.DB.Delete(&model.User{}, id).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败", err)
		return
	}
	response.OK(c, gin.H{"success": true})
}

func UpdateMe(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		response.Error(c, http.StatusBadRequest, "用户名不能为空", nil)
		return
	}

	err := db.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("username", req.Username).Error
	if err != nil {
		response.Error(c, http.StatusBadRequest, "用户名已存在", nil)
		return
	}

	token, _ := middleware.GenerateToken(user.ID, req.Username, user.Role)
	response.OK(c, gin.H{
		"token":    token,
		"username": req.Username,
	})
}
