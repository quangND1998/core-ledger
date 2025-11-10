package main

import (
	config "core-ledger/configs"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/handlers"
	"core-ledger/pkg/queue/jobs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load environment variables
	if err := os.Setenv("REDIS_ADDR", "localhost:6379"); err != nil {
		log.Fatal("Failed to set REDIS_ADDR")
	}
	if err := os.Setenv("REDIS_PASSWORD", ""); err != nil {
		log.Fatal("Failed to set REDIS_PASSWORD")
	}

	// Initialize Redis
	config.InitRedis()

	// Initialize Queue với error handling
	if err := queue.InitQueue(); err != nil {
		log.Fatalf("Failed to initialize queue: %v", err)
	}

	// Lấy config để khởi tạo worker
	queueConfig := config.GetQueueConfig()

	// Initialize Worker với pattern mới
	queue.InitWorker(queueConfig.RedisAddr, queueConfig.Concurrency, queueConfig.Queues)

	// Đăng ký các job handlers

	queue.RegisterJobHandler("data:process", &jobs.DataProcessJob{}, handlers.NewDataProcessHandler())

	// Start worker
	go func() {
		if err := queue.StartWorker(); err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()

	log.Println("Queue worker started successfully")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down queue worker...")
	queue.StopWorker()
	queue.CloseQueue()
	log.Println("Queue worker stopped")
}
