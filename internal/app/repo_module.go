package app

import (
	"core-ledger/pkg/repo"

	"go.uber.org/fx"
	// import thêm các repo khác
)

var RepoModule = fx.Module("repo",
	fx.Provide(
		repo.NewTransactionRepo,
	),
)
