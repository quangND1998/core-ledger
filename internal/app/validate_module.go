package app

import (
	"core-ledger/internal/module/validate"

	"go.uber.org/fx"
	// ... import thêm các service khác
)

var ValidateModule = fx.Module("validate",
	fx.Invoke(validate.ProvideValidator),
)
