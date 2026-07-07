package handler

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"panel/config"
	"panel/db"
	"panel/model"
	"panel/response"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件", nil)
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "icon"
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".webp": true, ".svg": true, ".ico": true,
	}
	if !allowedExts[ext] {
		response.Error(c, http.StatusBadRequest, "不支持的文件格式", nil)
		return
	}

	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano()/1e6, randomString(6), ext)
	dest := filepath.Join(config.UploadDir, filename)

	if err := c.SaveUploadedFile(file, dest); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存文件失败", err)
		return
	}

	image := model.Image{
		Filename:     filename,
		OriginalName: file.Filename,
		Category:     category,
	}
	err = db.DB.Create(&image).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "保存记录失败", err)
		return
	}

	response.OK(c, model.ImageItem{
		ID:           image.ID,
		Filename:     filename,
		OriginalName: file.Filename,
		Category:     category,
		URL:          "/uploads/" + filename,
	})
}

func GetImages(c *gin.Context) {
	category := c.Query("category")
	q := c.Query("q")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "50")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 200 {
		pageSize = 200
	}

	query := db.DB.Model(&model.Image{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if q != "" {
		query = query.Where("original_name LIKE ?", "%"+q+"%")
	}

	var total int64
	query.Count(&total)

	var images []model.Image
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&images).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	imageItems := make([]model.ImageItem, len(images))
	for i, img := range images {
		imageItems[i] = model.ImageItem{
			ID:           img.ID,
			Filename:     img.Filename,
			OriginalName: img.OriginalName,
			Category:     img.Category,
			URL:          "/uploads/" + img.Filename,
			CreatedAt:    img.CreatedAt,
		}
	}

	response.OK(c, model.ImageListResponse{
		Data:     imageItems,
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
	})
}

func UploadZip(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件", nil)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".zip" {
		response.Error(c, http.StatusBadRequest, "仅支持 .zip 压缩包", nil)
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "icon"
	}

	tempPath := filepath.Join(config.UploadDir, "temp_"+file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存文件失败", err)
		return
	}
	defer os.Remove(tempPath)

	r, err := zip.OpenReader(tempPath)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "解压失败: "+err.Error(), nil)
		return
	}
	defer r.Close()

	imgExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".webp": true, ".svg": true, ".ico": true,
	}

	var imported []model.ImageItem
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		entryExt := strings.ToLower(filepath.Ext(f.Name))
		if !imgExts[entryExt] {
			continue
		}

		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano()/1e6, randomString(6), entryExt)
		dest := filepath.Join(config.UploadDir, filename)

		rc, err := f.Open()
		if err != nil {
			continue
		}

		outFile, err := os.Create(dest)
		if err != nil {
			rc.Close()
			continue
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			continue
		}

		image := model.Image{
			Filename:     filename,
			OriginalName: f.Name,
			Category:     category,
		}
		if db.DB.Create(&image).Error != nil {
			continue
		}

		imported = append(imported, model.ImageItem{
			ID:           image.ID,
			Filename:     filename,
			OriginalName: f.Name,
			Category:     category,
			URL:          "/uploads/" + filename,
		})
	}

	response.OK(c, gin.H{
		"imported": len(imported),
		"images":   imported,
	})
}

func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	var image model.Image
	err := db.DB.First(&image, id).Error
	if err != nil {
		response.Error(c, http.StatusNotFound, "图片不存在", nil)
		return
	}

	filePath := filepath.Join(config.UploadDir, image.Filename)
	os.Remove(filePath)

	db.DB.Delete(&image)
	response.OK(c, gin.H{"success": true})
}

func ClearCategory(c *gin.Context) {
	category := c.Param("category")

	var images []model.Image
	err := db.DB.Where("category = ?", category).Find(&images).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "查询失败", err)
		return
	}

	for _, img := range images {
		os.Remove(filepath.Join(config.UploadDir, img.Filename))
	}

	result := db.DB.Where("category = ?", category).Delete(&model.Image{})
	response.OK(c, gin.H{
		"success": true,
		"deleted": result.RowsAffected,
	})
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
