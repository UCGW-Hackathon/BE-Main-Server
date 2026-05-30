package controllers

import (
	"net/http"
	"strconv"

	"whatsapp-backend/services"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
)

type HomeController interface {
	GetUserHome(ctx *gin.Context)
	GetWorkerHome(ctx *gin.Context)
}

type homeController struct {
	homeService services.HomeService
}

func NewHomeController(homeService services.HomeService) HomeController {
	return &homeController{homeService: homeService}
}

func (hc *homeController) GetUserHome(ctx *gin.Context) {
	lat, _ := strconv.ParseFloat(ctx.Query("latitude"), 64)
	lng, _ := strconv.ParseFloat(ctx.Query("longitude"), 64)
	data, err := hc.homeService.GetUserHome(ctx, lat, lng)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (hc *homeController) GetWorkerHome(ctx *gin.Context) {
	data, err := hc.homeService.GetWorkerHome(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}
