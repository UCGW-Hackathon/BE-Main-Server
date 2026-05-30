package services

import (
	"context"
	"time"

	"whatsapp-backend/models/entity"

	"gorm.io/gorm"
)

type NotificationService interface {
	ListNotifications(ctx context.Context) (any, error)
	MarkRead(ctx context.Context, notificationID string) error
	MarkAllRead(ctx context.Context) error
}

type notificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &notificationService{db: db}
}

func (s *notificationService) ListNotifications(ctx context.Context) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var notifications []entity.Notification
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(notifications))
	for _, notification := range notifications {
		data = append(data, map[string]any{
			"notification_id": notification.ID.String(),
			"type":            notification.Type,
			"title":           notification.Title,
			"body":            notification.Body,
			"deep_link":       notification.DeepLink,
			"metadata":        notification.Metadata,
			"is_read":         notification.IsRead,
			"read_at":         notification.ReadAt,
			"created_at":      notification.CreatedAt,
		})
	}
	return data, nil
}

func (s *notificationService) MarkRead(ctx context.Context, notificationID string) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}
	id, err := parseUUID(notificationID)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{"is_read": true, "read_at": &now}).Error
}

func (s *notificationService) MarkAllRead(ctx context.Context) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.db.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = FALSE", userID).
		Updates(map[string]any{"is_read": true, "read_at": &now}).Error
}
