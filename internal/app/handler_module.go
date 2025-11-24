package app

import (
	// "core-ledger/internal/auth/authhandler"
	// "core-ledger/internal/module/accounts/accounthandler"
	// "core-ledger/internal/module/transactions"
	// "core-ledger/internal/module/wallets"
	coaaccount "core-ledger/internal/module/coaAccount"
	"core-ledger/internal/module/entries"
	"core-ledger/internal/module/excel"
	"core-ledger/internal/module/permission"
	"core-ledger/internal/module/role"
	"core-ledger/internal/module/ruleCategory"
	"core-ledger/internal/module/ruleValue"
	"core-ledger/internal/module/swagger"
	"core-ledger/internal/module/transactions"
	"core-ledger/internal/module/user"

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
		// option.NewOptionHandler, // DEPRECATED: Không còn sử dụng model cũ
		permission.NewPermissionHandler,
		role.NewRoleHandler,
		user.NewUserHandler,
		user.NewAuthHandler,
		coaaccount.NewRequestCoaAccountHandler,
		swagger.NewSwaggerHandler,
	// accounthandler.NewAccountHandler,
	// authhandler.NewHandler,
	// wallets.NewWalletHandler,
	// transactions.NewTransactionHandler,
	),
)
