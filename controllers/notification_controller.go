package controllers

import (
	"net/http"

	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	ListNotifications(ctx *gin.Context)
	MarkRead(ctx *gin.Context)
	MarkAllRead(ctx *gin.Context)
}

type notificationController struct {
	notificationService services.NotificationService
}

func NewNotificationController(notificationService services.NotificationService) NotificationController {
	return &notificationController{notificationService: notificationService}
}

func (nc *notificationController) ListNotifications(ctx *gin.Context) {
	data, err := nc.notificationService.ListNotifications(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, map[string]any{
		"current_page": 1,
		"per_page":     20,
		"total":        0,
		"total_pages":  0,
		"unread_count": 0,
	})
}

func (nc *notificationController) MarkRead(ctx *gin.Context) {
	notificationID := ctx.Param("notification_id")
	err := nc.notificationService.MarkRead(ctx, notificationID)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Notifikasi ditandai sudah dibaca", gin.H{}, nil)
}

func (nc *notificationController) MarkAllRead(ctx *gin.Context) {
	err := nc.notificationService.MarkAllRead(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Semua notifikasi ditandai sudah dibaca", gin.H{}, nil)
}
