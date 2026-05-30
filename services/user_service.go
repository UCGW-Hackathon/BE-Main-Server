package services

import (
	"context"
	"time"

	"whatsapp-backend/dto"
	"whatsapp-backend/models/entity"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(ctx context.Context) (any, error)
	UpdateProfile(ctx context.Context, req dto.UpdateUserProfileRequest) (any, error)
	UpdateAvatar(ctx context.Context, avatarURL string) (any, error)
	UpdateLocation(ctx context.Context, req dto.UpdateUserLocationRequest) (any, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetProfile(ctx context.Context) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	var activeOrders int64
	_ = s.db.WithContext(ctx).Model(&entity.Order{}).
		Where("user_id = ? AND status IN ?", user.ID, []entity.OrderStatus{
			entity.OrderStatusPending,
			entity.OrderStatusAccepted,
			entity.OrderStatusOnTheWay,
			entity.OrderStatusArrived,
			entity.OrderStatusInProgress,
			entity.OrderStatusWorkPaused,
		}).
		Count(&activeOrders).Error

	return userResponse(user, map[string]any{
		"active_orders_count": activeOrders,
	}), nil
}

func (s *userService) UpdateProfile(ctx context.Context, req dto.UpdateUserProfileRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return map[string]any{
		"user_id":    user.ID.String(),
		"full_name":  user.FullName,
		"phone":      user.Phone,
		"address":    user.Address,
		"latitude":   user.Latitude,
		"longitude":  user.Longitude,
		"updated_at": user.UpdatedAt,
	}, nil
}

func (s *userService) UpdateAvatar(ctx context.Context, avatarURL string) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error; err != nil {
		return nil, err
	}
	return map[string]any{"avatar_url": avatarURL}, nil
}

func (s *userService) UpdateLocation(ctx context.Context, req dto.UpdateUserLocationRequest) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(map[string]any{
		"latitude":   req.Latitude,
		"longitude":  req.Longitude,
		"address":    req.Address,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}
	return map[string]any{
		"latitude":  req.Latitude,
		"longitude": req.Longitude,
		"address":   req.Address,
	}, nil
}

func userResponse(user entity.User, extra map[string]any) map[string]any {
	response := map[string]any{
		"user_id":           user.ID.String(),
		"full_name":         user.FullName,
		"email":             user.Email,
		"phone":             user.Phone,
		"role":              user.Role,
		"avatar_url":        user.AvatarURL,
		"address":           user.Address,
		"latitude":          user.Latitude,
		"longitude":         user.Longitude,
		"is_active":         user.IsActive,
		"email_verified_at": user.EmailVerifiedAt,
		"phone_verified_at": user.PhoneVerifiedAt,
		"last_login_at":     user.LastLoginAt,
		"created_at":        user.CreatedAt,
	}
	for key, value := range extra {
		response[key] = value
	}
	return response
}
