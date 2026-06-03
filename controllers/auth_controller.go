package controllers

import (
	"net/http"

	"situkang/dto"
	"situkang/services"
	"situkang/utils"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Refresh(ctx *gin.Context)
	Logout(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService: authService}
}

func (ac *authController) Register(ctx *gin.Context) {
	req := RequestJSON[dto.AuthRegisterRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := ac.authService.Register(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusCreated, "Registrasi berhasil", data, nil)
}

func (ac *authController) Login(ctx *gin.Context) {
	req := RequestJSON[dto.AuthLoginRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := ac.authService.Login(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (ac *authController) Refresh(ctx *gin.Context) {
	req := RequestJSON[dto.AuthRefreshRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	data, err := ac.authService.Refresh(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "", data, nil)
}

func (ac *authController) Logout(ctx *gin.Context) {
	err := ac.authService.Logout(ctx)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Berhasil logout", gin.H{}, nil)
}

func (ac *authController) ForgotPassword(ctx *gin.Context) {
	req := RequestJSON[dto.AuthForgotPasswordRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	err := ac.authService.ForgotPassword(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Link reset password telah dikirim", gin.H{}, nil)
}

func (ac *authController) ResetPassword(ctx *gin.Context) {
	req := RequestJSON[dto.AuthResetPasswordRequest](ctx)
	if ctx.IsAborted() {
		return
	}
	err := ac.authService.ResetPassword(ctx, req)
	if err != nil {
		utils.ResponseFAILED(ctx, nil, err)
		return
	}
	utils.JSONSuccess(ctx, http.StatusOK, "Password berhasil direset", gin.H{}, nil)
}
