package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                   uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderNumber          string       `gorm:"type:varchar(50);not null;uniqueIndex"`
	UserID               uuid.UUID    `gorm:"type:uuid;not null;index"`
	WorkerID             uuid.UUID    `gorm:"type:uuid;not null;index"`
	ServiceID            uuid.UUID    `gorm:"type:uuid;not null;index"`
	CategoryID           uuid.UUID    `gorm:"type:uuid;not null;index"`
	Title                string       `gorm:"type:varchar(255);not null"`
	Description          string       `gorm:"type:text;not null"`
	Status               OrderStatus  `gorm:"type:order_status;not null;default:'pending';index"`
	Urgency              OrderUrgency `gorm:"type:order_urgency;not null;default:'normal'"`
	LocationAddress      string       `gorm:"type:text;not null"`
	LocationDetail       *string      `gorm:"type:varchar(500)"`
	LocationLat          float64      `gorm:"type:decimal(10,7);not null"`
	LocationLng          float64      `gorm:"type:decimal(10,7);not null"`
	PreferredDate        *time.Time   `gorm:"type:date"`
	PreferredTimeStart   *time.Time   `gorm:"type:time"`
	PreferredTimeEnd     *time.Time   `gorm:"type:time"`
	Notes                *string      `gorm:"type:text"`
	BookingFee           int          `gorm:"type:integer;not null;default:2000"`
	BaseServiceFee       *int         `gorm:"type:integer"`
	TotalMaterialCost    int          `gorm:"type:integer;not null;default:0"`
	TotalAdditionalCost  int          `gorm:"type:integer;not null;default:0"`
	GrandTotal           *int         `gorm:"type:integer"`
	CancellationReason   *string      `gorm:"type:text"`
	CancellationCategory *string      `gorm:"type:varchar(50)"`
	CancelledBy          *uuid.UUID   `gorm:"type:uuid"`
	AcceptedAt           *time.Time   `gorm:"type:timestamptz"`
	StartedAt            *time.Time   `gorm:"type:timestamptz"`
	CompletedAt          *time.Time   `gorm:"type:timestamptz"`
	CancelledAt          *time.Time   `gorm:"type:timestamptz"`
	CreatedAt            time.Time    `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt            time.Time    `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type OrderPhoto struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID      uuid.UUID `gorm:"type:uuid;not null;index"`
	PhotoURL     string    `gorm:"type:text;not null"`
	Caption      *string   `gorm:"type:varchar(500)"`
	DisplayOrder int       `gorm:"type:integer;not null;default:0"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type OrderTimeline struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Event       string     `gorm:"type:varchar(50);not null"`
	Label       string     `gorm:"type:varchar(255);not null"`
	Description *string    `gorm:"type:text"`
	ActorID     *uuid.UUID `gorm:"type:uuid"`
	ActorType   *string    `gorm:"type:varchar(20)"`
	Metadata    JSONB      `gorm:"type:jsonb"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
}

func (OrderTimeline) TableName() string {
	return "order_timeline"
}

type ChatMessage struct {
	ID          uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID     uuid.UUID   `gorm:"type:uuid;not null;index"`
	SenderID    uuid.UUID   `gorm:"type:uuid;not null;index"`
	SenderType  string      `gorm:"type:varchar(20);not null"`
	Content     *string     `gorm:"type:text"`
	MessageType MessageType `gorm:"type:message_type;not null;default:'text'"`
	MediaURL    *string     `gorm:"type:text"`
	IsRead      bool        `gorm:"not null;default:false"`
	ReadAt      *time.Time  `gorm:"type:timestamptz"`
	CreatedAt   time.Time   `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_reviews_order_type"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	WorkerID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ReviewType string    `gorm:"type:varchar(20);not null;default:'worker';uniqueIndex:idx_reviews_order_type"`
	Rating     int16     `gorm:"type:smallint;not null;check:rating >= 1 AND rating <= 5"`
	Comment    *string   `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type ReviewTag struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ReviewID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Tag       string    `gorm:"type:varchar(50);not null;index"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type Invoice struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID              uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	InvoiceNumber        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	BaseServiceFee       int       `gorm:"type:integer;not null;default:0"`
	TotalMaterialCost    int       `gorm:"type:integer;not null;default:0"`
	TotalAdditionalCost  int       `gorm:"type:integer;not null;default:0"`
	BookingFee           int       `gorm:"type:integer;not null;default:2000"`
	PlatformFee          int       `gorm:"type:integer;not null;default:0"`
	DiscountAmount       int       `gorm:"type:integer;not null;default:0"`
	PromoCode            *string   `gorm:"type:varchar(50)"`
	GrandTotal           int       `gorm:"type:integer;not null;default:0"`
	Currency             string    `gorm:"type:varchar(3);not null;default:'IDR'"`
	PaymentInstruction   *string   `gorm:"type:text"`
	AIWorkSummary        *string   `gorm:"type:text"`
	AIMaterialsSummary   *string   `gorm:"type:text"`
	WorkerNotes          *string   `gorm:"type:text"`
	AllPurchasesApproved bool      `gorm:"not null;default:true"`
	CreatedAt            time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type InvoiceLineItem struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	InvoiceID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Label        string     `gorm:"type:varchar(255);not null"`
	Description  *string    `gorm:"type:text"`
	Category     string     `gorm:"type:varchar(50);not null"`
	Quantity     *float64   `gorm:"type:decimal(10,2);default:1"`
	Unit         *string    `gorm:"type:varchar(50)"`
	UnitPrice    *int       `gorm:"type:integer"`
	Amount       int        `gorm:"type:integer;not null;default:0"`
	PurchaseID   *uuid.UUID `gorm:"type:uuid;index"`
	DisplayOrder int        `gorm:"type:integer;not null;default:0"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;autoCreateTime"`
}

type Payment struct {
	ID              uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID         uuid.UUID     `gorm:"type:uuid;not null;index"`
	InvoiceID       uuid.UUID     `gorm:"type:uuid;not null;index"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index"`
	Amount          int           `gorm:"type:integer;not null"`
	Currency        string        `gorm:"type:varchar(3);not null;default:'IDR'"`
	PaymentMethod   PaymentMethod `gorm:"type:payment_method;not null"`
	PaymentStatus   PaymentStatus `gorm:"type:payment_status;not null;default:'unpaid';index"`
	PaymentProofURL *string       `gorm:"type:text"`
	TransactionRef  *string       `gorm:"type:varchar(255)"`
	SnapToken       *string       `gorm:"type:varchar(255)"`
	SnapRedirectURL *string       `gorm:"type:text"`
	PaidAt          *time.Time    `gorm:"type:timestamptz"`
	RefundedAt      *time.Time    `gorm:"type:timestamptz"`
	RefundAmount    *int          `gorm:"type:integer"`
	RefundReason    *string       `gorm:"type:text"`
	CreatedAt       time.Time     `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt       time.Time     `gorm:"type:timestamptz;not null;autoUpdateTime"`
}
