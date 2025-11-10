package app

import "go.uber.org/fx"

var FXProviders = fx.Options(
	CoreModule,
	RepoModule,
	ServiceModule,
	HandlerModule,
	RouterModule,
	JobModule,
	fx.Provide(NewApplication),
)
