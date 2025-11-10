package main

import (
	"context"
	config "core-ledger/configs"
	"core-ledger/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

func main() {
	// Load environment variables
	if err := os.Setenv("REDIS_ADDR", "localhost:6379"); err != nil {
		log.Fatal("Failed to set REDIS_ADDR")
	}
	if err := os.Setenv("REDIS_PASSWORD", ""); err != nil {
		log.Fatal("Failed to set REDIS_PASSWORD")
	}

	// Khởi tạo Redis (nếu các phần khác cần)
	config.InitRedis()

	// Dùng Fx để DI worker/handlers và auto start theo lifecycle
	fxApp := fx.New(
		app.CoreModule,
		app.RepoModule, // nếu cần tạo factory job cho nơi khác dùng
		app.QueueModule, // module worker + handler + lifecycle
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := fxApp.Start(ctx); err != nil {
		log.Fatalf("failed to start fx app: %v", err)
	}

	log.Println("Queue worker (Fx) started successfully")

	// Chờ tín hiệu hệ thống
	<-ctx.Done()

	log.Println("Shutting down queue worker (Fx)...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := fxApp.Stop(shutdownCtx); err != nil {
		log.Printf("error stopping fx app: %v", err)
	}
	log.Println("Queue worker stopped.")
}
