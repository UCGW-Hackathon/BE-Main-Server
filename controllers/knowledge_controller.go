package controllers

import (
	"net/http"

	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type KnowledgeController interface {
	ListArticles(ctx *gin.Context)
	GetArticle(ctx *gin.Context)
	ListFAQ(ctx *gin.Context)
}

type knowledgeController struct {
	knowledgeService services.KnowledgeService
}

func NewKnowledgeController(knowledgeService services.KnowledgeService) KnowledgeController {
	return &knowledgeController{knowledgeService: knowledgeService}
}

func (kc *knowledgeController) ListArticles(ctx *gin.Context) {
	data, err := kc.knowledgeService.ListArticles(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     10,
		"total":        0,
		"total_pages":  0,
	})
}

func (kc *knowledgeController) GetArticle(ctx *gin.Context) {
	articleID := ctx.Param("article_id")
	data, err := kc.knowledgeService.GetArticle(ctx, articleID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (kc *knowledgeController) ListFAQ(ctx *gin.Context) {
	data, err := kc.knowledgeService.ListFAQ(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}
