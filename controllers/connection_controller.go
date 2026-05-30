package controllers

import (
	"whatsapp-backend/dto"
	"whatsapp-backend/services"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
)

type ConnectionController interface {
	Connect(ctx *gin.Context)
}

type connectionController struct {
	connectionService services.ConnectionService
}

func NewConnectionController(connectionService services.ConnectionService) ConnectionController {
	return &connectionController{connectionService}
}

func (cc *connectionController) Connect(ctx *gin.Context) {
	var req dto.ConnectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendResponse[any, any](ctx, nil, nil, err)
		return
	}

	err := cc.connectionService.Connect(ctx, req)
	utils.SendResponse[any, any](ctx, nil, nil, err)
}
