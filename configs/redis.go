package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // ex: "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // "" nếu không có mật khẩu
		DB:       0,                           // default DB
	})

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")

}
