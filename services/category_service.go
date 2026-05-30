package services

import (
	"context"

	"whatsapp-backend/models/entity"

	"gorm.io/gorm"
)

type CategoryService interface {
	ListCategories(ctx context.Context) (any, error)
	ListCategoryServices(ctx context.Context, categoryID string) (any, error)
}

type categoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{db: db}
}

func (s *categoryService) ListCategories(ctx context.Context) (any, error) {
	var categories []entity.Category
	if err := s.db.WithContext(ctx).
		Where("is_active = TRUE").
		Order("display_order ASC, name ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(categories))
	for _, category := range categories {
		data = append(data, categoryResponse(category))
	}
	return data, nil
}

func (s *categoryService) ListCategoryServices(ctx context.Context, categoryID string) (any, error) {
	id, err := parseUUID(categoryID)
	if err != nil {
		return nil, err
	}

	var services []entity.Service
	if err := s.db.WithContext(ctx).
		Where("category_id = ? AND is_active = TRUE", id).
		Order("name ASC").
		Find(&services).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(services))
	for _, service := range services {
		data = append(data, serviceResponse(service))
	}
	return data, nil
}

func categoryResponse(category entity.Category) map[string]any {
	return map[string]any{
		"category_id":   category.ID.String(),
		"name":          category.Name,
		"slug":          category.Slug,
		"icon_url":      category.IconURL,
		"description":   category.Description,
		"display_order": category.DisplayOrder,
		"is_active":     category.IsActive,
	}
}

func serviceResponse(service entity.Service) map[string]any {
	return map[string]any{
		"service_id":         service.ID.String(),
		"category_id":        service.CategoryID.String(),
		"name":               service.Name,
		"slug":               service.Slug,
		"description":        service.Description,
		"icon_url":           service.IconURL,
		"base_price":         service.BasePrice,
		"price_unit":         service.PriceUnit,
		"estimated_duration": service.EstimatedDuration,
		"is_active":          service.IsActive,
	}
}
