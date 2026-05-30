package entity

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	ID                    uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID               uuid.UUID        `gorm:"type:uuid;not null;index"`
	WorkerID              uuid.UUID        `gorm:"type:uuid;not null;index"`
	ItemName              string           `gorm:"type:varchar(255);not null"`
	Category              PurchaseCategory `gorm:"type:purchase_category;not null;default:'material'"`
	Quantity              float64          `gorm:"type:decimal(10,2);not null;default:1"`
	Unit                  string           `gorm:"type:varchar(50);not null;default:'pcs'"`
	UnitPrice             int              `gorm:"type:integer;not null;default:0"`
	TotalPrice            int              `gorm:"type:integer;not null;default:0"`
	Reason                *string          `gorm:"type:text"`
	ReceiptPhotoURL       *string          `gorm:"type:text"`
	Status                PurchaseStatus   `gorm:"type:purchase_status;not null;default:'draft';index"`
	Confidence            *float64         `gorm:"type:decimal(3,2)"`
	NeedsClarification    bool             `gorm:"not null;default:false"`
	ClarificationQuestion *string          `gorm:"type:text"`
	ClarificationResponse *string          `gorm:"type:text"`
	AIExplanation         *string          `gorm:"type:text"`
	RawInput              *string          `gorm:"type:text"`
	AIProcessedAt         *time.Time       `gorm:"type:timestamptz"`
	ApprovedBy            *uuid.UUID       `gorm:"type:uuid"`
	ApprovedAt            *time.Time       `gorm:"type:timestamptz"`
	RejectedBy            *uuid.UUID       `gorm:"type:uuid"`
	RejectedAt            *time.Time       `gorm:"type:timestamptz"`
	RejectionReason       *string          `gorm:"type:text"`
	CreatedAt             time.Time        `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt             time.Time        `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type PurchaseRiskFlag struct {
	ID         uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchaseID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Type       RiskFlagType `gorm:"type:risk_flag_type;not null;index"`
	Message    string       `gorm:"type:text;not null"`
	IsResolved bool         `gorm:"not null;default:false"`
	ResolvedBy *uuid.UUID   `gorm:"type:uuid"`
	ResolvedAt *time.Time   `gorm:"type:timestamptz"`
	CreatedAt  time.Time    `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type PurchaseAuditLog struct {
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PurchaseID uuid.UUID   `gorm:"type:uuid;not null;index"`
	Action     AuditAction `gorm:"type:audit_action;not null;index"`
	ActorID    *uuid.UUID  `gorm:"type:uuid"`
	ActorName  *string     `gorm:"type:varchar(255)"`
	ActorType  string      `gorm:"type:varchar(20);not null"`
	Note       *string     `gorm:"type:text"`
	OldData    JSONB       `gorm:"type:jsonb"`
	NewData    JSONB       `gorm:"type:jsonb"`
	CreatedAt  time.Time   `gorm:"type:timestamptz;not null;autoCreateTime"`
}
