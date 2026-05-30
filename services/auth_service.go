package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"whatsapp-backend/config"
	"whatsapp-backend/dto"
	"whatsapp-backend/models/entity"
	http_error "whatsapp-backend/models/error"
	"whatsapp-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req dto.AuthRegisterRequest) (any, error)
	Login(ctx context.Context, req dto.AuthLoginRequest) (any, error)
	Refresh(ctx context.Context, req dto.AuthRefreshRequest) (any, error)
	Logout(ctx context.Context) error
	ForgotPassword(ctx context.Context, req dto.AuthForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.AuthResetPasswordRequest) error
}

type authService struct {
	db        *gorm.DB
	jwtConfig config.JWTConfig
	envConfig config.EnvConfig
}

func NewAuthService(db *gorm.DB, jwtConfig config.JWTConfig, envConfig config.EnvConfig) AuthService {
	return &authService{db: db, jwtConfig: jwtConfig, envConfig: envConfig}
}

func (s *authService) Register(ctx context.Context, req dto.AuthRegisterRequest) (any, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	var existing entity.User
	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, http_error.EMAIL_ALREADY_EXISTS
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Where("phone = ?", req.Phone).First(&existing).Error; err == nil {
		return nil, http_error.PHONE_ALREADY_EXISTS
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := entity.User{
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           req.Phone,
		PasswordHash:    string(passwordHash),
		Role:            entity.UserRole(req.Role),
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		IsActive:        true,
		EmailVerifiedAt: &now,
		PhoneVerifiedAt: &now,
	}

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if user.Role == entity.UserRoleWorker {
			profile := entity.WorkerProfile{
				UserID:             user.ID,
				VerificationStatus: entity.VerificationStatusUnverified,
				CertificateURLs:    entity.JSONB("[]"),
				IsAvailable:        true,
			}
			if err := tx.Create(&profile).Error; err != nil {
				return err
			}
			wallet := entity.WorkerWallet{
				WorkerID: user.ID,
				IsActive: true,
			}
			if err := tx.Create(&wallet).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.issueTokenResponse(ctx, user, true)
}

func (s *authService) Login(ctx context.Context, req dto.AuthLoginRequest) (any, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || strings.TrimSpace(req.Password) == "" {
		return nil, http_error.INVALID_CREDENTIALS
	}

	var user entity.User
	if err := s.db.WithContext(ctx).Where("email = ? AND is_active = TRUE", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http_error.INVALID_CREDENTIALS
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, http_error.INVALID_CREDENTIALS
	}

	now := time.Now()
	if err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Update("last_login_at", &now).Error; err != nil {
		return nil, err
	}
	user.LastLoginAt = &now

	return s.issueTokenResponse(ctx, user, true)
}

func (s *authService) Refresh(ctx context.Context, req dto.AuthRefreshRequest) (any, error) {
	tokenHash := hashToken(req.RefreshToken)
	var refresh entity.RefreshToken
	if err := s.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&refresh).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http_error.UNAUTHORIZED
		}
		return nil, err
	}

	var user entity.User
	if err := s.db.WithContext(ctx).Where("id = ? AND is_active = TRUE", refresh.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http_error.UNAUTHORIZED
		}
		return nil, err
	}

	now := time.Now()
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.RefreshToken{}).Where("id = ?", refresh.ID).Update("revoked_at", &now).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.issueTokenResponse(ctx, user, false)
}

func (s *authService) Logout(ctx context.Context) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.AuthForgotPasswordRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return http_error.VALIDATION_ERROR
	}
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.AuthResetPasswordRequest) error {
	if req.Token == "" || req.Password != req.PasswordConfirmation || !isStrongPassword(req.Password) {
		return http_error.VALIDATION_ERROR
	}
	return nil
}

func (s *authService) issueTokenResponse(ctx context.Context, user entity.User, includeUser bool) (map[string]any, error) {
	accessToken, expiresIn, err := utils.GenerateAccessToken(s.jwtConfig.GetSecretKey(), utils.TokenUser{
		ID:       user.ID,
		Role:     user.Role,
		Email:    user.Email,
		FullName: user.FullName,
	}, s.envConfig.GetAccessTokenTTL())
	if err != nil {
		return nil, err
	}

	refreshToken, err := randomToken()
	if err != nil {
		return nil, err
	}

	refresh := entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.envConfig.GetRefreshTokenTTL()),
	}

	if ginCtx, ok := ctx.(*gin.Context); ok {
		deviceInfo := ginCtx.GetHeader("User-Agent")
		ipAddress := ginCtx.ClientIP()
		if deviceInfo != "" {
			refresh.DeviceInfo = &deviceInfo
		}
		if ipAddress != "" {
			refresh.IPAddress = &ipAddress
		}
	}

	if err := s.db.WithContext(ctx).Create(&refresh).Error; err != nil {
		return nil, err
	}

	response := map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
	}
	if includeUser {
		response["user_id"] = user.ID.String()
		response["full_name"] = user.FullName
		response["email"] = user.Email
		response["phone"] = user.Phone
		response["role"] = user.Role
	}
	return response, nil
}

func validateRegisterRequest(req dto.AuthRegisterRequest) error {
	if req.FullName == "" {
		return http_error.VALIDATION_ERROR
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return http_error.VALIDATION_ERROR
	}
	if req.Phone == "" {
		return http_error.VALIDATION_ERROR
	}
	if req.Password != req.PasswordConfirmation {
		return http_error.VALIDATION_ERROR
	}
	if !isStrongPassword(req.Password) {
		return http_error.VALIDATION_ERROR
	}
	if req.Role != string(entity.UserRoleUser) && req.Role != string(entity.UserRoleWorker) {
		return http_error.VALIDATION_ERROR
	}
	return nil
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func randomToken() (string, error) {
	raw := make([]byte, 48)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func parseUUID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, http_error.BAD_REQUEST_ERROR
	}
	return id, nil
}
