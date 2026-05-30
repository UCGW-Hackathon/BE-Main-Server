package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null;index"`
	Type      NotificationType `gorm:"type:notification_type;not null;index"`
	Title     string           `gorm:"type:varchar(255);not null"`
	Body      string           `gorm:"type:text;not null"`
	DeepLink  *string          `gorm:"type:varchar(500)"`
	Metadata  JSONB            `gorm:"type:jsonb"`
	IsRead    bool             `gorm:"not null;default:false"`
	ReadAt    *time.Time       `gorm:"type:timestamptz"`
	CreatedAt time.Time        `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type Article struct {
	ID              uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title           string          `gorm:"type:varchar(500);not null"`
	Slug            string          `gorm:"type:varchar(500);not null;uniqueIndex"`
	Category        ArticleCategory `gorm:"type:article_category;not null;index"`
	ThumbnailURL    *string         `gorm:"type:text"`
	Excerpt         *string         `gorm:"type:text"`
	ContentHTML     string          `gorm:"type:text;not null"`
	ReadTimeMinutes *int            `gorm:"type:integer"`
	Author          *string         `gorm:"type:varchar(255);default:'Tim HandyDirect'"`
	Tags            JSONB           `gorm:"type:jsonb;default:'[]'"`
	IsPublished     bool            `gorm:"not null;default:false"`
	PublishedAt     *time.Time      `gorm:"type:timestamptz"`
	CreatedAt       time.Time       `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt       time.Time       `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type FAQ struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Question     string      `gorm:"type:text;not null"`
	Answer       string      `gorm:"type:text;not null"`
	Category     FAQCategory `gorm:"type:faq_category;not null;index"`
	DisplayOrder int         `gorm:"type:integer;not null;default:0"`
	IsActive     bool        `gorm:"not null;default:true"`
	CreatedAt    time.Time   `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt    time.Time   `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

func (FAQ) TableName() string {
	return "faqs"
}

type Promotion struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title           string     `gorm:"type:varchar(255);not null"`
	Description     *string    `gorm:"type:text"`
	ImageURL        string     `gorm:"type:text;not null"`
	CTALabel        *string    `gorm:"type:varchar(100)"`
	DeepLink        *string    `gorm:"type:varchar(500)"`
	PromoCode       *string    `gorm:"type:varchar(50)"`
	DiscountPercent *float64   `gorm:"type:decimal(5,2)"`
	DiscountAmount  *int       `gorm:"type:integer"`
	MinOrderAmount  *int       `gorm:"type:integer"`
	DisplayOrder    int        `gorm:"type:integer;not null;default:0"`
	IsActive        bool       `gorm:"not null;default:true"`
	ValidFrom       *time.Time `gorm:"type:timestamptz"`
	ValidUntil      *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"type:timestamptz;not null;autoUpdateTime"`
}
