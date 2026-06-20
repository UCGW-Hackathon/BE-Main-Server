package provider

import (
	"log"

	"situkang/config"
	"situkang/models/entity"
)

type ConfigProvider interface {
	ProvideJWTConfig() config.JWTConfig
	ProvideEnvConfig() config.EnvConfig
	ProvideDatabaseConfig() config.DatabaseConfig
}

type configProvider struct {
	jWTConfig      config.JWTConfig
	envConfig      config.EnvConfig
	databaseConfig config.DatabaseConfig
}

func NewConfigProvider() ConfigProvider {
	envConfig := config.NewEnvConfig("Asia/Jakarta")
	jWTConfig := config.NewJWTConfig(envConfig.GetJWTSecret())
	databaseConfig := config.NewDatabaseConfig(
		envConfig.GetDatabaseHost(),
		envConfig.GetDatabaseUser(),
		envConfig.GetDatabasePassword(),
		envConfig.GetDatabaseName(),
		envConfig.GetDatabasePort(),
		envConfig.GetDatabaseSSLMode())

	if envConfig.GetAppEnv() == "production" && len(envConfig.GetJWTSecret()) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters in production")
	}

	if envConfig.ShouldAutoMigrate() {
		if err := databaseConfig.AutoMigrateAll(
			&entity.User{},
			&entity.RefreshToken{},
			&entity.WorkerProfile{},
			&entity.Category{},
			&entity.Service{},
			&entity.WorkerService{},
			&entity.Order{},
			&entity.OrderPhoto{},
			&entity.OrderTimeline{},
			&entity.Purchase{},
			&entity.PurchaseRiskFlag{},
			&entity.PurchaseAuditLog{},
			&entity.ChatMessage{},
			&entity.Review{},
			&entity.ReviewTag{},
			&entity.Invoice{},
			&entity.InvoiceLineItem{},
			&entity.Payment{},
			&entity.Notification{},
			&entity.Article{},
			&entity.FAQ{},
			&entity.Promotion{},
			&entity.WorkerWallet{},
			&entity.WalletTransaction{},
			&entity.UploadedFile{},
		); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
	}

	if envConfig.ShouldSeedData() {
		if err := seedReferenceData(databaseConfig.GetInstance()); err != nil {
			log.Fatal("Failed to seed reference data:", err)
		}
	}

	return &configProvider{
		jWTConfig:      jWTConfig,
		envConfig:      envConfig,
		databaseConfig: databaseConfig,
	}
}

func (c *configProvider) ProvideJWTConfig() config.JWTConfig {
	return c.jWTConfig
}

func (c *configProvider) ProvideEnvConfig() config.EnvConfig {
	return c.envConfig
}

func (c *configProvider) ProvideDatabaseConfig() config.DatabaseConfig {
	return c.databaseConfig
}
