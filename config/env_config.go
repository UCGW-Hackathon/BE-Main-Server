package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type EnvConfig interface {
	GetAppEnv() string
	GetTCPAddress() string
	GetLogPath() string
	GetHostAddress() string
	GetHostPort() string
	GetEmailVerificationDuration() int
	GetDatabaseHost() string
	GetDatabasePort() string
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetDatabaseName() string
	GetDatabaseSSLMode() string
	GetSalt() string
	GetJWTSecret() string
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
	ShouldAutoMigrate() bool
	ShouldSeedData() bool
	GetUploadBaseURL() string
	GetSupabaseURL() string
	GetSupabaseKey() string
	GetSupabaseBucket() string
}

type envConfig struct {
	timezone string
}

func NewEnvConfig(timezone string) EnvConfig {
	godotenv.Load()
	os.Setenv("TZ", timezone)
	return &envConfig{
		timezone: timezone,
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil {
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes"
}

func (e *envConfig) GetAppEnv() string {
	return strings.ToLower(getEnv("APP_ENV", "development"))
}

func (e *envConfig) GetTCPAddress() string {
	host := e.GetHostAddress()
	port := e.GetHostPort()
	if host == "" {
		return ":" + port
	}
	return host + ":" + port
}

func (e *envConfig) GetLogPath() string {
	return getEnv("LOG_PATH", "logs")
}

func (e *envConfig) GetHostAddress() string {
	return strings.TrimSpace(os.Getenv("HOST_ADDRESS"))
}

func (e *envConfig) GetHostPort() string {
	return getEnv("HOST_PORT", "8080")
}

func (e *envConfig) GetEmailVerificationDuration() int {
	return getEnvInt("EMAIL_VERIFICATION_DURATION", 15)
}

func (e *envConfig) GetDatabaseHost() string {
	return getEnv("DB_HOST", "localhost")
}

func (e *envConfig) GetDatabasePort() string {
	return getEnv("DB_PORT", "5432")
}

func (e *envConfig) GetDatabaseUser() string {
	return getEnv("DB_USER", "postgres")
}

func (e *envConfig) GetDatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}

func (e *envConfig) GetDatabaseName() string {
	return getEnv("DB_NAME", "handydirect")
}

func (e *envConfig) GetDatabaseSSLMode() string {
	return getEnv("DB_SSLMODE", "disable")
}

func (e *envConfig) GetSalt() string {
	return getEnv("SALT", "")
}

func (e *envConfig) GetJWTSecret() string {
	secret := getEnv("JWT_SECRET", "")
	if secret != "" {
		return secret
	}
	return e.GetSalt()
}

func (e *envConfig) GetAccessTokenTTL() time.Duration {
	return time.Duration(getEnvInt("JWT_ACCESS_TOKEN_TTL_SECONDS", 3600)) * time.Second
}

func (e *envConfig) GetRefreshTokenTTL() time.Duration {
	return time.Duration(getEnvInt("JWT_REFRESH_TOKEN_TTL_HOURS", 720)) * time.Hour
}

func (e *envConfig) ShouldAutoMigrate() bool {
	return getEnvBool("AUTO_MIGRATE", true)
}

func (e *envConfig) ShouldSeedData() bool {
	return getEnvBool("SEED_DATA", true)
}

func (e *envConfig) GetUploadBaseURL() string {
	return strings.TrimRight(getEnv("UPLOAD_BASE_URL", "/uploads"), "/")
}

func (e *envConfig) GetSupabaseURL() string {
	return strings.TrimSpace(os.Getenv("SUPABASE_URL"))
}

func (e *envConfig) GetSupabaseKey() string {
	return strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_KEY"))
}

func (e *envConfig) GetSupabaseBucket() string {
	return strings.TrimSpace(os.Getenv("SUPABASE_BUCKET_NAME"))
}
