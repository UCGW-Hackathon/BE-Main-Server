package controllers

import (
	"math"
	"net/http"
	"strconv"

	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type FileController interface {
	Upload(ctx *gin.Context)
	GetFiles(ctx *gin.Context)
}

type fileController struct {
	fileService services.FileService
}

func NewFileController(fileService services.FileService) FileController {
	return &fileController{fileService: fileService}
}

func (c *fileController) Upload(ctx *gin.Context) {
	category := ctx.Param("category")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "File 'file' is required in multipart form data", nil)
		return
	}

	data, err := c.fileService.UploadFile(ctx.Request.Context(), ctx, fileHeader, category)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}

	utils.JSONSuccess(ctx, http.StatusCreated, "File berhasil diunggah", data, nil)
}

func (c *fileController) GetFiles(ctx *gin.Context) {
	category := ctx.Query("category")
	
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	data, total, err := c.fileService.GetFiles(ctx.Request.Context(), category, page, perPage)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	meta := map[string]any{
		"current_page": page,
		"per_page":     perPage,
		"total":        total,
		"total_pages":  totalPages,
	}

	utils.JSONSuccess(ctx, http.StatusOK, "File berhasil diambil", data, meta)
}
