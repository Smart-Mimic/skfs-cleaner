package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DryRun     bool
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	RouteID    string
	MaxCopies  int
}

func LoadEnv() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config := Config{
		DryRun:     getEnvAsBool("DRY_RUN", false),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUsername: getEnv("DB_USERNAME", "chirpstack"),
		DBPassword: getEnv("DB_PASSWORD", "chirpstack"),
		RouteID:    getEnv("ROUTE_ID", ""),
		MaxCopies:  getEnvAsInt("MAX_COPIES", 5),
	}

	if config.RouteID == "" {
		log.Fatal("ROUTE_ID is required but not set")
	}

	return config
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		// If it cannot be parsed, return the default value
		return defaultValue
	}
	return val
}
