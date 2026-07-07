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

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少用户名或密码", nil)
		return
	}

	var user model.User
	err := db.DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误", nil)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "生成令牌失败", err)
		return
	}

	response.OK(c, model.LoginResponse{
		Token:    token,
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
		Role:     user.Role,
	})
}

func ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "缺少旧密码或新密码", nil)
		return
	}

	if len(req.NewPassword) < 6 {
		response.Error(c, http.StatusBadRequest, "新密码至少 6 位", nil)
		return
	}

	user := middleware.GetUser(c)
	if user == nil {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	var dbUser model.User
	err := db.DB.Select("password").First(&dbUser, user.ID).Error
	if err != nil || bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.OldPassword)) != nil {
		response.Error(c, http.StatusBadRequest, "旧密码错误", nil)
		return
	}

	newHash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	err = db.DB.Model(&model.User{}).Where("id = ?", user.ID).Update("password", string(newHash)).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "修改密码失败", err)
		return
	}

	response.OK(c, gin.H{"success": true})
}
