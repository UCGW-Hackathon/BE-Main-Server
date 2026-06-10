package utils

import (
	"errors"
	"net/http"

	dto "situkang/models/dto"
	http_error "situkang/models/error"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func JSONSuccess[TData any](c *gin.Context, status int, message string, data TData, meta any) {
	c.JSON(status, dto.SuccessResponse[TData]{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func JSONError(c *gin.Context, status int, code string, message string, fieldErrors []dto.FieldError) {
	c.JSON(status, dto.ErrorResponse{
		Status:    "error",
		Message:   message,
		ErrorCode: code,
		Errors:    fieldErrors,
	})
}

func ResponseOK[Tdata any, TMetaData any](c *gin.Context, metaData TMetaData, data Tdata) {
	JSONSuccess(c, http.StatusOK, "", data, metaData)
}

func ResponseFAILED(c *gin.Context, metaData any, err error) {
	status, code, message := resolveError(err)
	JSONError(c, status, code, message, nil)
}

func SendResponse[Tdata any, TMetaData any](c *gin.Context, metaData TMetaData, data Tdata, err error) {
	if !c.IsAborted() {
		if err != nil {
			ResponseFAILED(c, metaData, err)
			c.Abort()
			return
		}
		ResponseOK(c, metaData, data)
		c.Abort()
		return
	}
}

func resolveError(err error) (int, string, string) {
	if errors.Is(err, http_error.BAD_REQUEST_ERROR) {
		return http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request format"
	}
	if errors.Is(err, http_error.VALIDATION_ERROR) {
		return http.StatusBadRequest, "VALIDATION_ERROR", err.Error()
	}
	if errors.Is(err, http_error.UNAUTHORIZED) {
		return http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized"
	}
	if errors.Is(err, http_error.INVALID_CREDENTIALS) || errors.Is(err, http_error.WRONG_PASSWORD) {
		return http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau password salah"
	}
	if errors.Is(err, http_error.FORBIDDEN) {
		return http.StatusForbidden, "FORBIDDEN", "Forbidden"
	}
	if errors.Is(err, http_error.DATA_NOT_FOUND) || errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, "NOT_FOUND", "Resource not found"
	}
	if errors.Is(err, http_error.EMAIL_ALREADY_EXISTS) {
		return http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email sudah terdaftar"
	}
	if errors.Is(err, http_error.PHONE_ALREADY_EXISTS) {
		return http.StatusConflict, "PHONE_ALREADY_EXISTS", "Nomor telepon sudah terdaftar"
	}
	if errors.Is(err, http_error.DUPLICATE_DATA) {
		return http.StatusConflict, "DUPLICATE_DATA", "Duplicate resource"
	}
	if errors.Is(err, http_error.INSUFFICIENT_BALANCE) {
		return http.StatusUnprocessableEntity, "INSUFFICIENT_BALANCE", "Saldo tidak mencukupi"
	}
	if errors.Is(err, http_error.INVALID_STATUS_TRANSITION) {
		return http.StatusUnprocessableEntity, "INVALID_STATUS_TRANSITION", "Transisi status tidak valid"
	}
	if errors.Is(err, http_error.CANCEL_NOT_ALLOWED) {
		return http.StatusUnprocessableEntity, "CANCEL_NOT_ALLOWED", "Pesanan tidak dapat dibatalkan karena worker sudah tiba di lokasi"
	}
	if errors.Is(err, http_error.TIMEOUT) {
		return http.StatusGatewayTimeout, "TIMEOUT", "Request timed out"
	}
	if errors.Is(err, http_error.NOT_IMPLEMENTED) {
		return http.StatusNotImplemented, "NOT_IMPLEMENTED", "Endpoint not implemented"
	}
	if errors.Is(err, http_error.INTERNAL_SERVER_ERROR) {
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	}
	return http.StatusInternalServerError, "INTERNAL_ERROR", err.Error()
}
