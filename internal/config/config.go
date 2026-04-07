package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	PostgresDSN         string
	RedisAddr           string
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

	// Format Data Source Name (DSN) untuk PostgreSQL
	postgresDSN := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	return &Config{
		AppPort:             os.Getenv("APP_PORT"),
		AppEnv:              os.Getenv("APP_ENV"),
		PostgresDSN:         postgresDSN,
		RedisAddr:           redisAddr,
		PasetoSecretKey:     os.Getenv("PASETO_SECRET_KEY"),
		XenditAPIKey:        os.Getenv("XENDIT_API_KEY"),
		XenditCallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
	}
}
