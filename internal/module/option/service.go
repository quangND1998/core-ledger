package option

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type OptionsSerive struct {
	db                 *gorm.DB
	logger             logger.CustomLogger
	ruleLayerRepo      repo.RuleLayerRepo
	ruleOptionStepRepo repo.RuleOptionStepRepo
	ruleOptionRepo     repo.RuleOptionRepo
	dispatcher         queue.Dispatcher
}

func NewOptionsSerive(dispatcher queue.Dispatcher, db *gorm.DB, ruleLayerRepo repo.RuleLayerRepo, ruleOptionStepRepo repo.RuleOptionStepRepo, ruleOptionRepo repo.RuleOptionRepo) *OptionsSerive {
	return &OptionsSerive{
		db:                 db,
		ruleLayerRepo:      ruleLayerRepo,
		ruleOptionStepRepo: ruleOptionStepRepo,
		ruleOptionRepo:     ruleOptionRepo,
		logger:             logger.NewSystemLog("OptionsSerive"),
		dispatcher:         dispatcher,
	}
}

func (s *OptionsSerive) GetRuleTypes(ctx context.Context) ([]*model.AccountRuleOption, error) {
	return s.ruleOptionRepo.GetManyByFields(ctx, map[string]any{
		"parent_option_id": nil,
		"status":           "ACTIVE",
	})

}

func (s *OptionsSerive) GetRuleGroups(ctx context.Context, typeId int64) ([]*model.AccountRuleOption, error) {
	return s.ruleOptionRepo.GetManyByFields(ctx, map[string]any{
		"parent_option_id": typeId,
		"status":           "ACTIVE",
	})

}
