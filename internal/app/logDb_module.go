package app

import (
	"core-ledger/internal/logging"

	"go.uber.org/fx"
	// ... import thêm các service khác
)

var LogDBModule = fx.Module("logDB_module",
	fx.Provide(
		logging.NewService,
	),
)
