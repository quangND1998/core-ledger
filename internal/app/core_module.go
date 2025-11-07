package app

import (
	"core-ledger/pkg/database"

	"go.uber.org/fx"
	// "core-ledger/internal/app/core/mail"
	// "core-ledger/internal/app/core/sqs"
)

var CoreModule = fx.Module("core",
	fx.Provide(
		database.Instance,
	),
)
