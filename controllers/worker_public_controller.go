package controllers

import (
	"net/http"
	"strconv"

	"whatsapp-backend/services"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
)

type WorkerPublicController interface {
	ListNearby(ctx *gin.Context)
	Search(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	GetReviews(ctx *gin.Context)
	GetServices(ctx *gin.Context)
}

type workerPublicController struct {
	workerService services.WorkerPublicService
}

func NewWorkerPublicController(workerService services.WorkerPublicService) WorkerPublicController {
	return &workerPublicController{workerService: workerService}
}

func (wc *workerPublicController) ListNearby(ctx *gin.Context) {
	lat, _ := strconv.ParseFloat(ctx.Query("latitude"), 64)
	lng, _ := strconv.ParseFloat(ctx.Query("longitude"), 64)
	data, err := wc.workerService.ListNearby(ctx, lat, lng)
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

func (wc *workerPublicController) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	lat, _ := strconv.ParseFloat(ctx.Query("latitude"), 64)
	lng, _ := strconv.ParseFloat(ctx.Query("longitude"), 64)
	data, err := wc.workerService.Search(ctx, query, lat, lng)
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

func (wc *workerPublicController) GetDetail(ctx *gin.Context) {
	workerID := ctx.Param("worker_id")
	data, err := wc.workerService.GetDetail(ctx, workerID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (wc *workerPublicController) GetReviews(ctx *gin.Context) {
	workerID := ctx.Param("worker_id")
	data, err := wc.workerService.GetReviews(ctx, workerID)
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

func (wc *workerPublicController) GetServices(ctx *gin.Context) {
	workerID := ctx.Param("worker_id")
	data, err := wc.workerService.GetServices(ctx, workerID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}
