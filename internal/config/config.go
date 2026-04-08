package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	DatabaseURL         string
	RedisURL            string
	PasetoSecretKey     string
	XenditAPIKey        string
	XenditCallbackToken string
}

func LoadConfig() *Config {
	// Load .env file jika ada (biasanya di environment lokal)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Config{
		AppPort:             os.Getenv("APP_PORT"),
		AppEnv:              os.Getenv("APP_ENV"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		RedisURL:            os.Getenv("REDIS_URL"),
		PasetoSecretKey:     os.Getenv("PASETO_SECRET_KEY"),
		XenditAPIKey:        os.Getenv("XENDIT_API_KEY"),
		XenditCallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
	}
}
