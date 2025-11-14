package app

import (
	coaaccount "core-ledger/internal/module/coaAccount"
	"core-ledger/internal/module/entries"
	"core-ledger/internal/module/excel"
	"core-ledger/internal/module/transactions"

	"go.uber.org/fx"
	// ... import thêm các service khác
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		transactions.NewTransactionService,
		excel.NewExcelService,
		coaaccount.NewCoaAccountService,
		entries.NewEntriesService,
	),
)
