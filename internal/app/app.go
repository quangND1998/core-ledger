package app

import (
	"context"
	config "core-ledger/configs"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"log"
	_ "time/tzdata"
)

type Application struct {
	log logger.CustomLogger
}

func NewApplication() *Application {
	return &Application{
		log: logger.NewSystemLog("Application"),
	}
}

func (a *Application) Run(ctx context.Context) error {
	a.log.Info("🚀 Application started...")
	config.InitRedis()
	if err := queue.InitQueue(); err != nil {
		log.Printf("❌ Error initializing queue: %v", err)
		return nil
	}

	// Ví dụ: gọi start HTTP server hoặc consumer
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.log.Info("🛑 Application stopped.")
	return nil
}
