package app

import (
	// "core-ledger/internal/auth/authhandler"
	// "core-ledger/internal/module/accounts/accounthandler"
	// "core-ledger/internal/module/transactions"
	// "core-ledger/internal/module/wallets"
	coaaccount "core-ledger/internal/module/coaAccount"
	"core-ledger/internal/module/entries"
	"core-ledger/internal/module/excel"
	"core-ledger/internal/module/ruleCategory"
	"core-ledger/internal/module/ruleValue"
	"core-ledger/internal/module/transactions"

	"go.uber.org/fx"
)

var HandlerModule = fx.Module("handler",
	fx.Provide(
		transactions.NewTransactionHandler,
		excel.NewExcelHandler,
		coaaccount.NewCoaAccountHandler,
		entries.NewEntriesHandler,
		ruleCategory.NewRuleCategoryHandler,
		ruleValue.NewRuleValueHandler,
	// accounthandler.NewAccountHandler,
	// authhandler.NewHandler,
	// wallets.NewWalletHandler,
	// transactions.NewTransactionHandler,
	),
)
