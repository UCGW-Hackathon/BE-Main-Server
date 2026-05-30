package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FullName        string         `gorm:"type:varchar(255);not null"`
	Email           string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	Phone           string         `gorm:"type:varchar(20);not null;uniqueIndex"`
	PasswordHash    string         `gorm:"type:varchar(255);not null"`
	Role            UserRole       `gorm:"type:user_role;not null;default:'user';index"`
	AvatarURL       *string        `gorm:"type:text"`
	Address         *string        `gorm:"type:text"`
	Latitude        *float64       `gorm:"type:decimal(10,7)"`
	Longitude       *float64       `gorm:"type:decimal(10,7)"`
	IsActive        bool           `gorm:"not null;default:true;index"`
	EmailVerifiedAt *time.Time     `gorm:"type:timestamptz"`
	PhoneVerifiedAt *time.Time     `gorm:"type:timestamptz"`
	LastLoginAt     *time.Time     `gorm:"type:timestamptz"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type RefreshToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenHash  string     `gorm:"type:varchar(255);not null;index"`
	DeviceInfo *string    `gorm:"type:varchar(500)"`
	IPAddress  *string    `gorm:"type:inet"`
	ExpiresAt  time.Time  `gorm:"type:timestamptz;not null;index"`
	RevokedAt  *time.Time `gorm:"type:timestamptz"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type WorkerProfile struct {
	ID                 uuid.UUID          `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID             uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex"`
	Specialization     *string            `gorm:"type:varchar(255)"`
	Bio                *string            `gorm:"type:text"`
	CoverPhotoURL      *string            `gorm:"type:text"`
	VerificationStatus VerificationStatus `gorm:"type:verification_status;not null;default:'unverified';index"`
	IDCardURL          *string            `gorm:"type:text"`
	CertificateURLs    JSONB              `gorm:"type:jsonb;default:'[]'"`
	BasePrice          *int               `gorm:"type:integer"`
	PriceUnit          *string            `gorm:"type:varchar(50);default:'per kunjungan'"`
	BookingFee         int                `gorm:"type:integer;not null;default:2000"`
	RatingAvg          float64            `gorm:"type:decimal(2,1);not null;default:0"`
	TotalReviews       int                `gorm:"type:integer;not null;default:0"`
	CompletedJobs      int                `gorm:"type:integer;not null;default:0"`
	IsAvailable        bool               `gorm:"not null;default:true;index"`
	VerifiedAt         *time.Time         `gorm:"type:timestamptz"`
	CreatedAt          time.Time          `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt          time.Time          `gorm:"type:timestamptz;not null;autoUpdateTime"`
}
