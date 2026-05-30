package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name         string    `gorm:"type:varchar(100);not null"`
	Slug         string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	IconURL      *string   `gorm:"type:text"`
	Description  *string   `gorm:"type:text"`
	DisplayOrder int       `gorm:"type:integer;not null;default:0"`
	IsActive     bool      `gorm:"not null;default:true;index"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type Service struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Name              string    `gorm:"type:varchar(255);not null"`
	Slug              string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Description       *string   `gorm:"type:text"`
	IconURL           *string   `gorm:"type:text"`
	BasePrice         *int      `gorm:"type:integer"`
	PriceUnit         *string   `gorm:"type:varchar(50);default:'per kunjungan'"`
	EstimatedDuration *string   `gorm:"type:varchar(50)"`
	IsActive          bool      `gorm:"not null;default:true;index"`
	CreatedAt         time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type WorkerService struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkerID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_worker_services_worker_service"`
	ServiceID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_worker_services_worker_service"`
	CustomPrice *int      `gorm:"type:integer"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}
