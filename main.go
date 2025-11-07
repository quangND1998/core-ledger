package main

import (
	"context"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"core-ledger/internal/app"
)

func main() {
	fx.New(
		app.FXProviders,     // gom tất cả module Provider vào
		fx.Invoke(StartApp), // lifecycle của app
	).Run()
}

func StartApp(lc fx.Lifecycle, application *app.Application) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return application.Run(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return application.Shutdown(ctx)
		},
	})
}
