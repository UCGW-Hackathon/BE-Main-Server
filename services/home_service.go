package services

import (
	"context"
	"strings"

	"whatsapp-backend/models/entity"

	"gorm.io/gorm"
)

type HomeService interface {
	GetUserHome(ctx context.Context, latitude float64, longitude float64) (any, error)
	GetWorkerHome(ctx context.Context) (any, error)
}

type homeService struct {
	db                  *gorm.DB
	categoryService     CategoryService
	workerPublicService WorkerPublicService
}

func NewHomeService(db *gorm.DB, categoryService CategoryService, workerPublicService WorkerPublicService) HomeService {
	return &homeService{db: db, categoryService: categoryService, workerPublicService: workerPublicService}
}

func (s *homeService) GetUserHome(ctx context.Context, latitude float64, longitude float64) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	var activeOrder entity.Order
	activeOrderData := any(nil)
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND status IN ?", userID, []entity.OrderStatus{
			entity.OrderStatusPending,
			entity.OrderStatusAccepted,
			entity.OrderStatusOnTheWay,
			entity.OrderStatusArrived,
			entity.OrderStatusInProgress,
			entity.OrderStatusWorkPaused,
		}).
		Order("created_at DESC").
		First(&activeOrder).Error; err == nil {
		var worker entity.User
		var service entity.Service
		_ = s.db.WithContext(ctx).First(&worker, "id = ?", activeOrder.WorkerID).Error
		_ = s.db.WithContext(ctx).First(&service, "id = ?", activeOrder.ServiceID).Error
		activeOrderData = map[string]any{
			"order_id":     activeOrder.ID.String(),
			"status":       activeOrder.Status,
			"worker_name":  worker.FullName,
			"service_name": service.Name,
			"eta_minutes":  nil,
		}
	}

	var promotions []entity.Promotion
	if err := s.db.WithContext(ctx).
		Where("is_active = TRUE").
		Order("display_order ASC, created_at DESC").
		Limit(5).
		Find(&promotions).Error; err != nil {
		return nil, err
	}
	promoData := make([]map[string]any, 0, len(promotions))
	for _, promo := range promotions {
		promoData = append(promoData, map[string]any{
			"promo_id":    promo.ID.String(),
			"title":       promo.Title,
			"description": promo.Description,
			"image_url":   promo.ImageURL,
			"cta_label":   promo.CTALabel,
			"deep_link":   promo.DeepLink,
			"valid_until": promo.ValidUntil,
		})
	}

	var articles []entity.Article
	if err := s.db.WithContext(ctx).
		Where("is_published = TRUE").
		Order("published_at DESC NULLS LAST, created_at DESC").
		Limit(5).
		Find(&articles).Error; err != nil {
		return nil, err
	}
	articleData := make([]map[string]any, 0, len(articles))
	for _, article := range articles {
		articleData = append(articleData, map[string]any{
			"article_id":    article.ID.String(),
			"title":         article.Title,
			"thumbnail_url": article.ThumbnailURL,
			"cta_label":     "Baca",
			"slug":          article.Slug,
		})
	}

	categories, err := s.categoryService.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	workers, err := s.workerPublicService.ListNearby(ctx, latitude, longitude)
	if err != nil {
		return nil, err
	}

	displayName := user.FullName
	if parts := strings.Fields(user.FullName); len(parts) > 0 {
		displayName = parts[0]
	}

	return map[string]any{
		"user_summary": map[string]any{
			"full_name":       displayName,
			"avatar_url":      user.AvatarURL,
			"current_address": user.Address,
		},
		"active_order":       activeOrderData,
		"promotions":         promoData,
		"articles":           articleData,
		"service_categories": categories,
		"featured_workers":   workers,
	}, nil
}

func (s *homeService) GetWorkerHome(ctx context.Context) (any, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := s.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	var profile entity.WorkerProfile
	_ = s.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error

	var incomingCount int64
	_ = s.db.WithContext(ctx).Model(&entity.Order{}).
		Where("worker_id = ? AND status = ?", userID, entity.OrderStatusPending).
		Count(&incomingCount).Error

	return map[string]any{
		"worker_summary": map[string]any{
			"full_name":      user.FullName,
			"avatar_url":     user.AvatarURL,
			"is_available":   profile.IsAvailable,
			"rating":         profile.RatingAvg,
			"completed_jobs": profile.CompletedJobs,
		},
		"incoming_orders_count": incomingCount,
	}, nil
}
