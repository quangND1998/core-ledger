package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := RedisClient.Ping(ctx).Err(); err == nil {
			log.Println("Connected to Redis successfully")
			return
		}
		log.Printf("[WARN] Failed to connect Redis (attempt %d/5), retrying...", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Printf("[ERROR] Could not connect to Redis after retries, continuing without Redis")
	RedisClient = nil

}
