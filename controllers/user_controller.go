package controllers

import (
	"net/http"

	"situkang/dto"
	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetMe(ctx *gin.Context)
	UpdateMe(ctx *gin.Context)
	UpdateAvatar(ctx *gin.Context)
	UpdateLocation(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (uc *userController) GetMe(ctx *gin.Context) {
	data, err := uc.userService.GetProfile(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (uc *userController) UpdateMe(ctx *gin.Context) {
	req := RequestJSON[dto.UpdateUserProfileRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := uc.userService.UpdateProfile(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Profil berhasil diperbarui", data, nil)
}

func (uc *userController) UpdateAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Avatar file is required", nil)
		return
	}

	avatarURL, err := utils.SaveUploadedFile(ctx, file, "avatars")
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}
	data, err := uc.userService.UpdateAvatar(ctx, avatarURL)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Foto profil berhasil diperbarui", data, nil)
}

func (uc *userController) UpdateLocation(ctx *gin.Context) {
	req := RequestJSON[dto.UpdateUserLocationRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := uc.userService.UpdateLocation(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Lokasi berhasil diperbarui", data, nil)
}
