package middleware

import (
	"os"
	"strings"

	"situkang/models/entity"
	"situkang/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UserIDContextKey   = "user_id"
	RoleContextKey     = "role"
	EmailContextKey    = "email"
	FullNameContextKey = "full_name"
)

var jwtSecret string

func ConfigureAuth(secret string) {
	jwtSecret = strings.TrimSpace(secret)
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Missing Authorization header", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == authHeader {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Authorization header must use Bearer token", nil)
			c.Abort()
			return
		}

		secret := jwtSecret
		if secret == "" {
			secret = strings.TrimSpace(os.Getenv("JWT_SECRET"))
			if secret == "" {
				secret = strings.TrimSpace(os.Getenv("SALT"))
			}
		}

		claims, err := utils.ParseAccessToken(secret, tokenString)
		if err != nil {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Token tidak valid atau sudah expired", nil)
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Token tidak valid atau sudah expired", nil)
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Set(RoleContextKey, claims.Role)
		c.Set(EmailContextKey, claims.Email)
		c.Set(FullNameContextKey, claims.FullName)
		c.Next()
	}
}

func RequireRoles(roles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, ok := c.Get(RoleContextKey)
		if !ok {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Unauthorized", nil)
			c.Abort()
			return
		}

		roleStr, ok := roleValue.(string)
		if !ok {
			utils.JSONError(c, 401, "UNAUTHORIZED", "Unauthorized", nil)
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if roleStr == string(allowed) || roleStr == string(entity.UserRoleAdmin) {
				c.Next()
				return
			}
		}

		utils.JSONError(c, 403, "FORBIDDEN", "Forbidden", nil)
		c.Abort()
	}
}
