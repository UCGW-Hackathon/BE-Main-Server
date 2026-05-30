package services

import (
	"context"
	"math"
	"sort"
	"strings"

	"whatsapp-backend/models/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerPublicService interface {
	ListNearby(ctx context.Context, latitude float64, longitude float64) (any, error)
	Search(ctx context.Context, query string, latitude float64, longitude float64) (any, error)
	GetDetail(ctx context.Context, workerID string) (any, error)
	GetReviews(ctx context.Context, workerID string) (any, error)
	GetServices(ctx context.Context, workerID string) (any, error)
}

type workerPublicService struct {
	db *gorm.DB
}

func NewWorkerPublicService(db *gorm.DB) WorkerPublicService {
	return &workerPublicService{db: db}
}

func (s *workerPublicService) ListNearby(ctx context.Context, latitude float64, longitude float64) (any, error) {
	return s.listWorkers(ctx, "", latitude, longitude)
}

func (s *workerPublicService) Search(ctx context.Context, query string, latitude float64, longitude float64) (any, error) {
	return s.listWorkers(ctx, query, latitude, longitude)
}

func (s *workerPublicService) GetDetail(ctx context.Context, workerID string) (any, error) {
	id, err := parseUUID(workerID)
	if err != nil {
		return nil, err
	}

	row, err := s.getWorkerRow(ctx, id)
	if err != nil {
		return nil, err
	}

	data := workerResponse(row, 0)
	services, err := s.GetServices(ctx, workerID)
	if err != nil {
		return nil, err
	}
	reviews, err := s.GetReviews(ctx, workerID)
	if err != nil {
		return nil, err
	}
	data["services"] = services
	data["recent_reviews"] = reviews
	data["bio"] = row.Bio
	data["cover_photo_url"] = row.CoverPhotoURL
	data["booking_fee"] = row.BookingFee
	return data, nil
}

func (s *workerPublicService) GetReviews(ctx context.Context, workerID string) (any, error) {
	id, err := parseUUID(workerID)
	if err != nil {
		return nil, err
	}

	type reviewRow struct {
		ReviewID  uuid.UUID
		OrderID   uuid.UUID
		UserID    uuid.UUID
		FullName  string
		AvatarURL *string
		Rating    int16
		Comment   *string
		CreatedAt string
	}

	var rows []reviewRow
	if err := s.db.WithContext(ctx).Table("reviews r").
		Select("r.id as review_id, r.order_id, r.user_id, u.full_name, u.avatar_url, r.rating, r.comment, r.created_at").
		Joins("JOIN users u ON u.id = r.user_id").
		Where("r.worker_id = ?", id).
		Order("r.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		data = append(data, map[string]any{
			"review_id":  row.ReviewID.String(),
			"order_id":   row.OrderID.String(),
			"user_id":    row.UserID.String(),
			"full_name":  row.FullName,
			"avatar_url": row.AvatarURL,
			"rating":     row.Rating,
			"comment":    row.Comment,
			"created_at": row.CreatedAt,
		})
	}
	return data, nil
}

func (s *workerPublicService) GetServices(ctx context.Context, workerID string) (any, error) {
	id, err := parseUUID(workerID)
	if err != nil {
		return nil, err
	}

	type serviceRow struct {
		ServiceID         uuid.UUID
		CategoryID        uuid.UUID
		Name              string
		Slug              string
		Description       *string
		IconURL           *string
		BasePrice         *int
		CustomPrice       *int
		PriceUnit         *string
		EstimatedDuration *string
	}

	var rows []serviceRow
	if err := s.db.WithContext(ctx).Table("worker_services ws").
		Select("s.id as service_id, s.category_id, s.name, s.slug, s.description, s.icon_url, s.base_price, ws.custom_price, s.price_unit, s.estimated_duration").
		Joins("JOIN services s ON s.id = ws.service_id").
		Where("ws.worker_id = ? AND ws.is_active = TRUE AND s.is_active = TRUE", id).
		Order("s.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		price := row.BasePrice
		if row.CustomPrice != nil {
			price = row.CustomPrice
		}
		data = append(data, map[string]any{
			"service_id":         row.ServiceID.String(),
			"category_id":        row.CategoryID.String(),
			"name":               row.Name,
			"slug":               row.Slug,
			"description":        row.Description,
			"icon_url":           row.IconURL,
			"base_price":         price,
			"price_unit":         row.PriceUnit,
			"estimated_duration": row.EstimatedDuration,
		})
	}
	return data, nil
}

func (s *workerPublicService) listWorkers(ctx context.Context, query string, latitude float64, longitude float64) (any, error) {
	var rows []workerRow
	db := s.db.WithContext(ctx).Table("users u").
		Select(`u.id as worker_id, u.full_name, u.avatar_url, u.latitude, u.longitude,
			wp.specialization, wp.bio, wp.cover_photo_url, wp.verification_status,
			wp.base_price, wp.price_unit, wp.booking_fee, wp.rating_avg, wp.total_reviews,
			wp.completed_jobs, wp.is_available`).
		Joins("JOIN worker_profiles wp ON wp.user_id = u.id").
		Where("u.role = ? AND u.is_active = TRUE", entity.UserRoleWorker)

	query = strings.TrimSpace(strings.ToLower(query))
	if query != "" {
		like := "%" + query + "%"
		db = db.Where(`LOWER(u.full_name) LIKE ? OR LOWER(COALESCE(wp.specialization, '')) LIKE ? OR EXISTS (
			SELECT 1 FROM worker_services ws
			JOIN services s ON s.id = ws.service_id
			WHERE ws.worker_id = u.id AND LOWER(s.name) LIKE ?
		)`, like, like, like)
	}

	if err := db.Scan(&rows).Error; err != nil {
		return nil, err
	}

	data := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		distance := distanceKm(latitude, longitude, row.Latitude, row.Longitude)
		item := workerResponse(row, distance)
		services, _ := s.GetServices(ctx, row.WorkerID.String())
		item["services"] = serviceNames(services)
		data = append(data, item)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i]["distance_km"].(float64) < data[j]["distance_km"].(float64)
	})

	return data, nil
}

func (s *workerPublicService) getWorkerRow(ctx context.Context, id uuid.UUID) (workerRow, error) {
	var row workerRow
	err := s.db.WithContext(ctx).Table("users u").
		Select(`u.id as worker_id, u.full_name, u.avatar_url, u.latitude, u.longitude,
			wp.specialization, wp.bio, wp.cover_photo_url, wp.verification_status,
			wp.base_price, wp.price_unit, wp.booking_fee, wp.rating_avg, wp.total_reviews,
			wp.completed_jobs, wp.is_available`).
		Joins("JOIN worker_profiles wp ON wp.user_id = u.id").
		Where("u.id = ? AND u.role = ? AND u.is_active = TRUE", id, entity.UserRoleWorker).
		First(&row).Error
	return row, err
}

type workerRow struct {
	WorkerID           uuid.UUID
	FullName           string
	AvatarURL          *string
	Latitude           *float64
	Longitude          *float64
	Specialization     *string
	Bio                *string
	CoverPhotoURL      *string
	VerificationStatus entity.VerificationStatus
	BasePrice          *int
	PriceUnit          *string
	BookingFee         int
	RatingAvg          float64
	TotalReviews       int
	CompletedJobs      int
	IsAvailable        bool
}

func workerResponse(row workerRow, distance float64) map[string]any {
	return map[string]any{
		"worker_id":      row.WorkerID.String(),
		"full_name":      row.FullName,
		"avatar_url":     row.AvatarURL,
		"specialization": row.Specialization,
		"rating":         row.RatingAvg,
		"total_reviews":  row.TotalReviews,
		"completed_jobs": row.CompletedJobs,
		"distance_km":    math.Round(distance*10) / 10,
		"is_verified":    row.VerificationStatus == entity.VerificationStatusVerified,
		"is_available":   row.IsAvailable,
		"base_price":     row.BasePrice,
		"price_unit":     row.PriceUnit,
		"latitude":       row.Latitude,
		"longitude":      row.Longitude,
	}
}

func distanceKm(lat1 float64, lng1 float64, lat2 *float64, lng2 *float64) float64 {
	if lat2 == nil || lng2 == nil || lat1 == 0 || lng1 == 0 {
		return 0
	}

	const earthRadiusKm = 6371
	dLat := degToRad(*lat2 - lat1)
	dLng := degToRad(*lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(*lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func degToRad(value float64) float64 {
	return value * math.Pi / 180
}

func serviceNames(services any) []string {
	items, ok := services.([]map[string]any)
	if !ok {
		return []string{}
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if name, ok := item["name"].(string); ok {
			names = append(names, name)
		}
	}
	return names
}
