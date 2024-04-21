package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBName       string
	DBPort       int
	DBUser       string
	DBPass       string
	APIKey       string
	Port         int
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	SourceEmail  string
	TargetEmail  string
}

// LoadConfig reads environment variables and returns a config object
func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
		return nil, err
	}

	config := &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnvAsInt("DB_PORT", 27017),
		Port:         getEnvAsInt("PORT", 3000),
		DBName:       os.Getenv("DB_NAME"),
		DBUser:       os.Getenv("DB_USER"),
		DBPass:       os.Getenv("DB_PASS"),
		APIKey:       os.Getenv("API_KEY"),
		AWSRegion:    os.Getenv("AWS_REGION"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretKey: os.Getenv("AWS_SECRET_KEY"),
		SourceEmail:  os.Getenv("SOURCE_EMAIL"),
		TargetEmail:  os.Getenv("TARGET_EMAIL"),
	}

	return config, nil
}

// Helper functions to read and parse environment variables
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
