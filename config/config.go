package config

import (
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	Port           int
	IsProduction   bool
	AllowedOrigins []string
	MaxRoomSize    int
}

var cfg *AppConfig

func Get() *AppConfig {
	return cfg
}

func Load() {
	log.Println("Loading config from environment...")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
	}

	isProduction := os.Getenv("IS_PRODUCTION") == "false"

	maxRoomSize, err := strconv.Atoi(os.Getenv("MAX_ROOM_SIZE"))
	if err != nil || maxRoomSize <= 0 {
		maxRoomSize = 2
	}

	cfg = &AppConfig{
		Port:           port,
		IsProduction:   isProduction,
		AllowedOrigins: []string{"*"},
		MaxRoomSize:    maxRoomSize,
	}
}
