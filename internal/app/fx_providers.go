package app

import "go.uber.org/fx"

var FXProviders = fx.Options(
	CoreModule,
	RepoModule,
	ServiceModule,
	HandlerModule,
	RouterModule,
	QueueClientModule,
	ValidateModule,
	LogDBModule,
	// QueueModule,
	fx.Provide(NewApplication),
)
