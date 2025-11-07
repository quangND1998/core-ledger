package app

import (
	"core-ledger/internal/module/transactions"

	"go.uber.org/fx"
	// ... import thêm các service khác
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		transactions.NewTransactionService,
	),
)
