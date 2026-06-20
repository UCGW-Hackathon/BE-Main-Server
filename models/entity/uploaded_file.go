package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadedFile struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Filename    string         `gorm:"type:varchar(255);not null" json:"filename"`
	StoredName  string         `gorm:"type:varchar(255);not null" json:"stored_name"`
	Path        string         `gorm:"type:text;not null" json:"path"`
	URL         string         `gorm:"type:text;not null" json:"url"`
	Category    string         `gorm:"type:varchar(50);not null;index" json:"category"` // "cover_photo", "avatar", "damage_proof", "invoice"
	Size        int64          `gorm:"type:bigint;not null" json:"size"`
	ContentType string         `gorm:"type:varchar(100);not null" json:"content_type"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
