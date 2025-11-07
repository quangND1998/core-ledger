package app

import (
	// "core-ledger/internal/auth/authhandler"
	// "core-ledger/internal/module/accounts/accounthandler"
	// "core-ledger/internal/module/transactions"
	// "core-ledger/internal/module/wallets"
	"core-ledger/internal/module/transactions"

	"go.uber.org/fx"
)

var HandlerModule = fx.Module("handler",
	fx.Provide(
		transactions.NewTransactionHandler,
	// accounthandler.NewAccountHandler,
	// authhandler.NewHandler,
	// wallets.NewWalletHandler,
	// transactions.NewTransactionHandler,
	),
)
