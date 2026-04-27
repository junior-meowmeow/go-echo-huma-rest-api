package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	Auth   AuthConfig
	Log    LogConfig
	Mongo  MongoConfig
	S3     S3Config
	Client ClientConfig
}

type AppConfig struct {
	Port               int    `env:"PORT" envDefault:"8000"`
	APIBasePath        string `env:"API_BASE_PATH"`
	RequestPerSecLimit int    `env:"REQ_PER_SEC_LIMIT" envDefault:"20"`
	EnablePrometheus   bool   `env:"ENABLE_PROMETHEUS" envDefault:"true"`
}

type AuthConfig struct {
	JWTSecret       string        `env:"JWT_SECRET" envDefault:"test-secret"`
	TokenExpiration time.Duration `env:"TOKEN_EXPIRATION" envDefault:"72h"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type MongoConfig struct {
	Host   string `env:"MONGO_HOST" envDefault:"mongo"`
	Port   string `env:"MONGO_PORT" envDefault:"27017"`
	DBName string `env:"DB_NAME" envDefault:"testdb"`
	DBUser string `env:"DB_USER" envDefault:"user"`
	DBPass string `env:"DB_PASS" envDefault:"pass"`
}

type S3Config struct {
	Endpoint string `env:"S3_ENDPOINT" envDefault:"http://localhost:8333"`
	Bucket   string `env:"S3_BUCKET" envDefault:"test-bucket"`
}

type ClientConfig struct {
	PetStoreURL string `env:"PETSTORE_URL" envDefault:"http://localhost:8080/api/v3"`
}

// NewConfig loads environment variables and returns a Config.
func NewConfig() (Config, error) {
	envType := os.Getenv("APP_ENV")
	if envType == "" {
		envType = "local"
	}

	if envType == "local" {
		basePath := "./config"
		// Load env files with precedence.
		if err := godotenv.Load(filepath.Join(basePath, ".env.local")); err != nil {
			slog.Warn("Could not load .env.local", slog.Any("error", err))
		}
		if err := godotenv.Load(filepath.Join(basePath, ".env")); err != nil {
			slog.Warn("Could not load .env", slog.Any("error", err))
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}
