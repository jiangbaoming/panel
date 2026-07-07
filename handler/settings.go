package handler

import (
	"net/http"
	"panel/db"
	"panel/model"
	"panel/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))

	var settings model.UserSettings
	err := db.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		newSettings := model.UserSettings{
			UserID:      userID,
			BgImage:     "",
			DisplayMode: "both",
			PageTitle:   "个人导航页",
			PageFavicon: "",
			Footer:      "",
		}
		db.DB.Create(&newSettings)
		settings = newSettings
	}

	response.OK(c, settings)
}

func UpdateSettings(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userId"))

	var req struct {
		BgImage        string `json:"bg_image"`
		DisplayMode    string `json:"display_mode"`
		PageTitle      string `json:"page_title"`
		PageFavicon    string `json:"page_favicon"`
		Footer         string `json:"footer"`
		WelcomeMessage string `json:"welcome_message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var settings model.UserSettings
	err := db.DB.Where("user_id = ?", userID).First(&settings).Error

	if err != nil {
		settings = model.UserSettings{
			UserID:         userID,
			BgImage:        req.BgImage,
			DisplayMode:    req.DisplayMode,
			PageTitle:      req.PageTitle,
			PageFavicon:    req.PageFavicon,
			Footer:         req.Footer,
			WelcomeMessage: req.WelcomeMessage,
		}
		db.DB.Create(&settings)
	} else {
		updates := map[string]interface{}{}
		if req.BgImage != "" {
			updates["bg_image"] = req.BgImage
		}
		if req.DisplayMode != "" {
			updates["display_mode"] = req.DisplayMode
		}
		if req.PageTitle != "" {
			updates["page_title"] = req.PageTitle
		}
		if req.PageFavicon != "" {
			updates["page_favicon"] = req.PageFavicon
		}
		if req.Footer != "" {
			updates["footer"] = req.Footer
		}
		if req.WelcomeMessage != "" {
			updates["welcome_message"] = req.WelcomeMessage
		}
		if len(updates) > 0 {
			db.DB.Model(&model.UserSettings{}).Where("user_id = ?", userID).Updates(updates)
		}
	}

	db.DB.Where("user_id = ?", userID).First(&settings)
	response.OK(c, settings)
}
