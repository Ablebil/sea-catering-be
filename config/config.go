package config

import (
	"log"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV"`
	AppHost string `env:"APP_HOST"`
	AppPort int    `env:"APP_PORT"`
	AppUrl  string `env:"APP_URL"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`

	FEURL string `env:"FE_URL"`

	AccessSecret  string `env:"ACCESS_SECRET"`
	RefreshSecret string `env:"REFRESH_SECRET"`

	EmailUser     string `env:"EMAIL_USER"`
	EmailPassword string `env:"EMAIL_PASSWORD"`

	OTPExpiry time.Duration `env:"OTP_EXPIRY"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`
	FERedirectURL      string `env:"FE_REDIRECT_URL"`

	StateLength int           `env:"STATE_LENGTH"`
	StateExpiry time.Duration `env:"STATE_EXPIRY"`

	SupabaseURL string `env:"SUPABASE_URL"`
	SupabaseKey string `env:"SUPABASE_KEY"`

	MaxFileSize int `env:"MAX_FILE_SIZE"`

	MidtransClientKey       string        `env:"MIDTRANS_CLIENT_KEY"`
	MidtransServerKey       string        `env:"MIDTRANS_SERVER_KEY"`
	MidtransPaymentDuration time.Duration `env:"MIDTRANS_PAYMENT_DURATION"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	log.Printf("DB Config - Host: %s, Port: %d, Name: %s, User: %s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser)

	return cfg, nil
}
