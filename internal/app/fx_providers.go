package app

import "go.uber.org/fx"

var FXProviders = fx.Options(
	CoreModule,
	RepoModule,
	ServiceModule,
	HandlerModule,
	RouterModule,
	fx.Provide(NewApplication),
)
