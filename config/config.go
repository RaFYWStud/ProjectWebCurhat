package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port           int
	IsProduction   bool
	AllowedOrigins []string
	MaxRoomSize    int
	DbURI          string
	JWTSecret      string
	AccessTokenTTL int64 // in seconds
}

var cfg *AppConfig

func Get() *AppConfig {
	return cfg
}

func Load() {
	log.Println("Loading config from environment...")

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file, using OS environment variables: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	isProduction := os.Getenv("IS_PRODUCTION") == "true"

	maxRoomSize, err := strconv.Atoi(os.Getenv("MAX_ROOM_SIZE"))
	if err != nil || maxRoomSize <= 0 {
		maxRoomSize = 2
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me"
		log.Println("[WARN] JWT_SECRET not set, using insecure default for development")
	}

	accessTokenTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil || accessTokenTTL <= 0 {
		accessTokenTTL = 86400 // 24 hours
	}

	cfg = &AppConfig{
		Port:           port,
		IsProduction:   isProduction,
		AllowedOrigins: []string{"*"},
		MaxRoomSize:    maxRoomSize,
		DbURI:          loadDatabaseConfig(),
		JWTSecret:      jwtSecret,
		AccessTokenTTL: int64(accessTokenTTL),
	}
}

func loadDatabaseConfig() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	user := getEnvOrDefault("DB_USER", "postgres")
	pass := getEnvOrDefault("DB_PASS", "postgres")
	name := getEnvOrDefault("DB_NAME", "webcurhat")
	port := getEnvOrDefault("DB_PORT", "5432")
	timeZone := getEnvOrDefault("DB_TIME_ZONE", "Asia/Makassar")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=disable",
		host, user, pass, name, port, timeZone)
}

func getEnvOrDefault(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
