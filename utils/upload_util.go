package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, subdir string) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is required")
	}
	if file.Size > 5*1024*1024 {
		return "", fmt.Errorf("file exceeds 5MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".pdf" {
		return "", fmt.Errorf("unsupported file type")
	}

	cleanSubdir := strings.Trim(filepath.Clean(subdir), ".\\/")
	targetDir := filepath.Join("uploads", cleanSubdir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.NewString(), ext)
	targetPath := filepath.Join(targetDir, name)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(targetPath), nil
}
