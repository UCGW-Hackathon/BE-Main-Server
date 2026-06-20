package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"situkang/config"
	"situkang/dto"
	"situkang/models/entity"
	http_error "situkang/models/error"
	"situkang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileService interface {
	UploadFile(ctx context.Context, ginCtx *gin.Context, fileHeader *multipart.FileHeader, category string) (*dto.FileResponse, error)
	GetFiles(ctx context.Context, category string, page, perPage int) ([]dto.FileResponse, int64, error)
}

type fileService struct {
	db  *gorm.DB
	env config.EnvConfig
}

func NewFileService(db *gorm.DB, env config.EnvConfig) FileService {
	return &fileService{db: db, env: env}
}

func (s *fileService) UploadFile(ctx context.Context, ginCtx *gin.Context, fileHeader *multipart.FileHeader, category string) (*dto.FileResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Validate category and map to subdirectory
	var subdir string
	switch category {
	case "cover_photo":
		subdir = "covers"
	case "avatar":
		subdir = "avatars"
	case "damage_proof":
		subdir = "damage_proofs"
	case "invoice":
		subdir = "invoices"
	default:
		return nil, fmt.Errorf("%w: invalid file category, must be one of: cover_photo, avatar, damage_proof, invoice", http_error.VALIDATION_ERROR)
	}

	// Validate extensions strictly
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if category != "invoice" && ext == ".pdf" {
		return nil, fmt.Errorf("%w: pdf file is only allowed for invoices", http_error.VALIDATION_ERROR)
	}

	// Save file locally using the upload utility
	relativePath, err := utils.SaveUploadedFile(ginCtx, fileHeader, subdir)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", http_error.VALIDATION_ERROR, err)
	}

	// Construct full URL
	appURL := s.env.GetAppURL()
	fullURL := appURL + relativePath

	// Save to DB
	// Stored name is the end segment of relativePath
	storedName := filepath.Base(relativePath)

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uploaded := entity.UploadedFile{
		UserID:      userID,
		Filename:    fileHeader.Filename,
		StoredName:  storedName,
		Path:        relativePath,
		URL:         fullURL,
		Category:    category,
		Size:        fileHeader.Size,
		ContentType: contentType,
	}

	if err := s.db.WithContext(ctx).Create(&uploaded).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", http_error.INTERNAL_SERVER_ERROR, err)
	}

	return &dto.FileResponse{
		ID:          uploaded.ID.String(),
		UserID:      uploaded.UserID.String(),
		Filename:    uploaded.Filename,
		URL:         uploaded.URL,
		Category:    uploaded.Category,
		Size:        uploaded.Size,
		ContentType: uploaded.ContentType,
		CreatedAt:   uploaded.CreatedAt,
	}, nil
}

func (s *fileService) GetFiles(ctx context.Context, category string, page, perPage int) ([]dto.FileResponse, int64, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	query := s.db.WithContext(ctx).Model(&entity.UploadedFile{}).Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var files []entity.UploadedFile
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Limit(perPage).Offset(offset).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	responses := make([]dto.FileResponse, 0, len(files))
	for _, f := range files {
		responses = append(responses, dto.FileResponse{
			ID:          f.ID.String(),
			UserID:      f.UserID.String(),
			Filename:    f.Filename,
			URL:         f.URL,
			Category:    f.Category,
			Size:        f.Size,
			ContentType: f.ContentType,
			CreatedAt:   f.CreatedAt,
		})
	}

	return responses, total, nil
}
