package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig interface {
	EnsureSchema() error
	AutoMigrateAll(entities ...interface{}) error
	GetInstance() *gorm.DB
}
type databaseConfig struct {
	db *gorm.DB
}

func NewDatabaseConfig(DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE string) DatabaseConfig {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta ",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	db = db.Session(&gorm.Session{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database pool:", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return &databaseConfig{db: db}
}

func (cfg *databaseConfig) EnsureSchema() error {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
		createEnumSQL("user_role", []string{"user", "worker", "admin"}),
		createEnumSQL("order_status", []string{"pending", "accepted", "on_the_way", "arrived", "in_progress", "work_paused", "completed", "cancelled", "rejected", "waiting_payment", "waiting_for_payment"}),
		createEnumSQL("order_urgency", []string{"normal", "urgent"}),
		createEnumSQL("purchase_status", []string{"draft", "pending_approval", "approved", "rejected", "needs_clarification"}),
		createEnumSQL("purchase_category", []string{"material", "alat", "sparepart", "bahan_bangunan", "biaya_tambahan", "lainnya"}),
		createEnumSQL("risk_flag_type", []string{"harga_tidak_wajar", "item_tidak_relevan", "data_tidak_lengkap", "nota_tidak_jelas", "duplikat", "alasan_tidak_lengkap"}),
		createEnumSQL("message_type", []string{"text", "image", "system"}),
		createEnumSQL("payment_method", []string{"cash", "bank_transfer", "ewallet"}),
		createEnumSQL("payment_status", []string{"unpaid", "pending", "paid", "refunded", "waiting_payment", "waiting_for_payment"}),
		createEnumSQL("verification_status", []string{"unverified", "pending", "verified", "rejected"}),
		createEnumSQL("notification_type", []string{"order", "purchase", "chat", "promo", "system", "payment"}),
		createEnumSQL("article_category", []string{"faq", "guide", "tips", "safety", "payment"}),
		createEnumSQL("faq_category", []string{"general", "payment", "tracking", "security", "cancellation"}),
		createEnumSQL("wallet_tx_type", []string{"earning", "withdrawal", "refund", "bonus", "fee"}),
		createEnumSQL("wallet_tx_status", []string{"pending", "completed", "failed", "cancelled"}),
		createEnumSQL("audit_action", []string{"created", "ai_processed", "submitted", "approved", "rejected", "clarification_requested", "clarification_responded", "edited", "deleted"}),
		`ALTER TYPE order_status ADD VALUE IF NOT EXISTS 'waiting_payment'`,
		`ALTER TYPE order_status ADD VALUE IF NOT EXISTS 'waiting_for_payment'`,
		`ALTER TYPE payment_status ADD VALUE IF NOT EXISTS 'waiting_payment'`,
		`ALTER TYPE payment_status ADD VALUE IF NOT EXISTS 'waiting_for_payment'`,
	}

	for _, statement := range statements {
		if err := cfg.db.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}

func (cfg *databaseConfig) AutoMigrateAll(entities ...interface{}) error {
	if err := cfg.EnsureSchema(); err != nil {
		return err
	}

	err := cfg.db.AutoMigrate(
		entities...,
	)

	return err

}

func (cfg *databaseConfig) GetInstance() *gorm.DB {
	return cfg.db
}

func createEnumSQL(name string, values []string) string {
	quotedValues := ""
	for idx, value := range values {
		if idx > 0 {
			quotedValues += ", "
		}
		quotedValues += fmt.Sprintf("'%s'", value)
	}

	return fmt.Sprintf(`
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
		CREATE TYPE %s AS ENUM (%s);
	END IF;
END
$$;`, name, name, quotedValues)
}
