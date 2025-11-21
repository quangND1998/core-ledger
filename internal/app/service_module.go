package app

import (
	coaaccount "core-ledger/internal/module/coaAccount"
	"core-ledger/internal/module/entries"
	"core-ledger/internal/module/excel"
	"core-ledger/internal/module/permission"
	"core-ledger/internal/module/role"
	"core-ledger/internal/module/ruleCategory"
	"core-ledger/internal/module/ruleValue"
	"core-ledger/internal/module/transactions"
	"core-ledger/internal/module/user"

	"go.uber.org/fx"
	// ... import thêm các service khác
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		transactions.NewTransactionService,
		excel.NewExcelService,
		coaaccount.NewCoaAccountService,
		entries.NewEntriesService,
		ruleCategory.NewRuleCateogySerive,
		ruleValue.NewRuleCateogySerive,
		// option.NewOptionsSerive, // DEPRECATED: Không còn sử dụng model cũ
		permission.NewPermissionService,
		role.NewRoleService,
		user.NewUserService,
		coaaccount.NewRequestCoaAccountService,
		coaaccount.NewRuleValidationService,
	),
)
