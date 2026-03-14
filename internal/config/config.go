package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	APIBasePath string
	MongoHost   string
	MongoPort   string
	MongoUser   string
	MongoPass   string
	DBName      string
	DBUser      string
	DBPass      string
	S3Endpoint  string
	S3Bucket    string
	PetStoreURL string
}

// Load environment variables and returns a Config. (Simple version)
func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	portStr := getEnv("PORT", "8888")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid PORT: %v", err)
	}

	return Config{
		Port:        port,
		APIBasePath: getEnv("API_BASE_PATH", ""),
		MongoHost:   getEnv("MONGO_HOST", "mongo"),
		MongoPort:   getEnv("MONGO_PORT", "27017"),
		MongoUser:   getEnv("MONGO_USER", "user"),
		MongoPass:   getEnv("MONGO_PASS", "pass"),
		DBName:      getEnv("DB_NAME", "test"),
		DBUser:      getEnv("DB_USER", "user"),
		DBPass:      getEnv("DB_PASS", "pass"),
		S3Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:8333"),
		S3Bucket:    getEnv("S3_BUCKET", "test-bucket"),
		PetStoreURL: getEnv("PETSTORE_URL", "http://localhost:8080/api/v3"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
