package services

import (
	"context"

	"whatsapp-backend/middleware"
	http_error "whatsapp-backend/models/error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func currentUserID(ctx context.Context) (uuid.UUID, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return uuid.Nil, http_error.UNAUTHORIZED
	}

	value, ok := ginCtx.Get(middleware.UserIDContextKey)
	if !ok {
		return uuid.Nil, http_error.UNAUTHORIZED
	}

	userID, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil, http_error.UNAUTHORIZED
	}

	return userID, nil
}

func currentRole(ctx context.Context) (string, error) {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return "", http_error.UNAUTHORIZED
	}

	value, ok := ginCtx.Get(middleware.RoleContextKey)
	if !ok {
		return "", http_error.UNAUTHORIZED
	}

	role, ok := value.(string)
	if !ok {
		return "", http_error.UNAUTHORIZED
	}

	return role, nil
}
