package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) *redis.Client {
	// ParseURL otomatis mendeteksi 'rediss://' dan menyalakan TLS/SSL Configuration
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Gagal memilah URL Redis: %v", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return rdb
}
