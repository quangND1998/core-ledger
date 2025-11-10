package config

import (
	"fmt"
	"os"
	"strconv"
)

// QueueConfig chứa cấu hình cho queue system
type QueueConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
	EnableWorker  bool
	Queues        map[string]int
}

// GetQueueConfig trả về cấu hình queue từ environment variables
func GetQueueConfig() *QueueConfig {
	redisAddr := Reader().Get("REDIS_ADDR")
	redisPassword := Reader().Get("REDIS_PASSWORD")
	redisDB := getEnvAsInt("REDIS_DB", 0)
	concurrency := getEnvAsInt("QUEUE_CONCURRENCY", 10)
	enableWorker := getEnvAsBool("QUEUE_WORKER_ENABLED", true)

	// Cấu hình queues với priority
	queues := map[string]int{
		"critical": getEnvAsInt("QUEUE_CRITICAL_WORKERS", 6),
		"default":  getEnvAsInt("QUEUE_DEFAULT_WORKERS", 3),
		"low":      getEnvAsInt("QUEUE_LOW_WORKERS", 1),
	}

	return &QueueConfig{
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		Concurrency:   concurrency,
		EnableWorker:  enableWorker,
		Queues:        queues,
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// ValidateQueueConfig kiểm tra cấu hình queue có hợp lệ không
func ValidateQueueConfig(config *QueueConfig) error {
	if config.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}

	if config.Concurrency <= 0 {
		return fmt.Errorf("QUEUE_CONCURRENCY must be greater than 0")
	}

	// Kiểm tra tổng số workers không vượt quá concurrency
	totalWorkers := 0
	for _, workers := range config.Queues {
		totalWorkers += workers
	}

	if totalWorkers > config.Concurrency {
		return fmt.Errorf("total queue workers (%d) cannot exceed concurrency (%d)", totalWorkers, config.Concurrency)
	}

	return nil
}

// GetQueueConfigWithValidation trả về cấu hình queue đã được validate
func GetQueueConfigWithValidation() (*QueueConfig, error) {
	config := GetQueueConfig()
	if err := ValidateQueueConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}
