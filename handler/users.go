package handler

import (
	"net/http"
	"panel/db"
	"panel/middleware"
	"panel/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 获取所有用户（管理员）
func GetUsers(c *gin.Context) {
	var users []model.User
	err := db.DB.Select("id, username, avatar, role").Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// 添加用户（管理员）
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Avatar   string `json:"avatar"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少用户名或密码"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": req.Username,
		"avatar":   req.Avatar,
		"role":     req.Role,
	})
}

// 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := db.DB.Delete(&model.User{}, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 修改自己的用户名
func UpdateMe(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	err := db.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("username", req.Username).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 签发新 token
	token, _ := middleware.GenerateToken(user.ID, req.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": req.Username,
	})
}
