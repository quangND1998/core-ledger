package app

import (
	"core-ledger/pkg/repo"

	"go.uber.org/fx"
	// import thêm các repo khác
)

var RepoModule = fx.Module("repo",
	fx.Provide(
		repo.NewTransactionRepo,
		repo.NewUserRepoIml,
		repo.NewCoAccountRepo,
		repo.NewEnTriesRepo,
		repo.NewTransactionLogRepo,
		repo.NewSnapshotRepo,
		repo.NewJournalRepo,
		repo.NewRuleCategoryRepo,
		repo.NewRuleValueRepo,
		repo.NewRuleOptionRepo,
		repo.NewRuleOptionStepRepo,
		repo.NewRuleLayerRepo,
		repo.NewPermissionRepo,
		repo.NewRoleRepo,
	),
)
