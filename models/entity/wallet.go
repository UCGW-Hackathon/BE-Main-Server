package entity

import (
	"time"

	"github.com/google/uuid"
)

type WorkerWallet struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkerID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Balance         int64     `gorm:"type:bigint;not null;default:0"`
	TotalEarnings   int64     `gorm:"type:bigint;not null;default:0"`
	TotalWithdrawn  int64     `gorm:"type:bigint;not null;default:0"`
	PendingEarnings int64     `gorm:"type:bigint;not null;default:0"`
	IsActive        bool      `gorm:"not null;default:true"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;not null;autoUpdateTime"`
}

type WalletTransaction struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	OrderID       *uuid.UUID     `gorm:"type:uuid;index"`
	Type          WalletTxType   `gorm:"type:wallet_tx_type;not null;index"`
	Amount        int            `gorm:"type:integer;not null"`
	BalanceBefore int64          `gorm:"type:bigint;not null"`
	BalanceAfter  int64          `gorm:"type:bigint;not null"`
	Description   *string        `gorm:"type:text"`
	ReferenceID   *string        `gorm:"type:varchar(255)"`
	Status        WalletTxStatus `gorm:"type:wallet_tx_status;not null;default:'pending';index"`
	CompletedAt   *time.Time     `gorm:"type:timestamptz"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;autoUpdateTime"`
}
