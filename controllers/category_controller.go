package controllers

import (
	"net/http"

	"whatsapp-backend/services"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	ListCategories(ctx *gin.Context)
	ListCategoryServices(ctx *gin.Context)
}

type categoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) CategoryController {
	return &categoryController{categoryService: categoryService}
}

func (cc *categoryController) ListCategories(ctx *gin.Context) {
	data, err := cc.categoryService.ListCategories(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     20,
		"total":        0,
		"total_pages":  0,
	})
}

func (cc *categoryController) ListCategoryServices(ctx *gin.Context) {
	categoryID := ctx.Param("category_id")
	data, err := cc.categoryService.ListCategoryServices(ctx, categoryID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}
