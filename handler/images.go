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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传图片
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "icon"
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".webp": true, ".svg": true, ".ico": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式"})
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano()/1e6, randomString(6), ext)
	dest := filepath.Join(config.UploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 写入数据库
	image := model.Image{
		Filename:     filename,
		OriginalName: file.Filename,
		Category:     category,
	}
	err = db.DB.Create(&image).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存记录失败"})
		return
	}

	c.JSON(http.StatusOK, model.ImageItem{
		ID:           image.ID,
		Filename:     filename,
		OriginalName: file.Filename,
		Category:     category,
		URL:          "/uploads/" + filename,
	})
}

// 获取图片列表
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

	// 构建查询
	query := db.DB.Model(&model.Image{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if q != "" {
		query = query.Where("original_name LIKE ?", "%"+q+"%")
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询数据
	var images []model.Image
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&images).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 转换为响应格式
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

	c.JSON(http.StatusOK, model.ImageListResponse{
		Data:     imageItems,
		Total:    int(total),
		Page:     page,
		PageSize: pageSize,
	})
}

// 上传 ZIP 压缩包
func UploadZip(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持 .zip 压缩包"})
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "icon"
	}

	// 保存临时文件
	tempPath := filepath.Join(config.UploadDir, "temp_"+file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer os.Remove(tempPath)

	// 解压
	r, err := zip.OpenReader(tempPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解压失败: " + err.Error()})
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

		// 生成唯一文件名
		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano()/1e6, randomString(6), entryExt)
		dest := filepath.Join(config.UploadDir, filename)

		// 读取并保存
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

		// 写入数据库
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

	c.JSON(http.StatusOK, gin.H{
		"imported": len(imported),
		"images":   imported,
	})
}

// 删除图片
func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	var image model.Image
	err := db.DB.First(&image, id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图片不存在"})
		return
	}

	// 删除文件
	filePath := filepath.Join(config.UploadDir, image.Filename)
	os.Remove(filePath)

	// 删除记录
	db.DB.Delete(&image)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 清空分类
func ClearCategory(c *gin.Context) {
	category := c.Param("category")

	// 查询所有该分类的图片
	var images []model.Image
	err := db.DB.Where("category = ?", category).Find(&images).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 删除文件
	for _, img := range images {
		os.Remove(filepath.Join(config.UploadDir, img.Filename))
	}

	// 删除记录
	result := db.DB.Where("category = ?", category).Delete(&model.Image{})
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deleted": result.RowsAffected,
	})
}

// 生成随机字符串
func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
