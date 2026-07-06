package config

import (
	"os"
	"path/filepath"
)

var (
	Port      string
	JWTSecret string
	DataDir   string
	DBPath    string
	UploadDir string
)

func Init() {
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "5678"
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "dev-jwt-secret-change-in-production"
	}

	DataDir = os.Getenv("DATA_DIR")
	if DataDir == "" {
		DataDir = "data"
	}

	DBPath = filepath.Join(DataDir, "data.db")

	UploadDir = os.Getenv("UPLOAD_DIR")
	if UploadDir == "" {
		UploadDir = filepath.Join(DataDir, "uploads")
	}

	// 确保目录存在
	os.MkdirAll(DataDir, 0755)
	os.MkdirAll(UploadDir, 0755)
}
