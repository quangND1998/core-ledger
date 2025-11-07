package app

import (
	"context"
	"core-ledger/pkg/logger"
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
	// Ví dụ: gọi start HTTP server hoặc consumer
	return nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.log.Info("🛑 Application stopped.")
	return nil
}
