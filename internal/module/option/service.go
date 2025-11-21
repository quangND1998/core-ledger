package option

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"

	"gorm.io/gorm"
)

type OptionsSerive struct {
	db         *gorm.DB
	logger     logger.CustomLogger
	dispatcher queue.Dispatcher
}

// NewOptionsSerive - DEPRECATED: Service này đã không còn sử dụng model cũ
// Có thể xóa hoàn toàn nếu không cần thiết
func NewOptionsSerive(dispatcher queue.Dispatcher, db *gorm.DB) *OptionsSerive {
	return &OptionsSerive{
		db:         db,
		logger:     logger.NewSystemLog("OptionsSerive"),
		dispatcher: dispatcher,
	}
}

// GetRuleTypes - DEPRECATED: Sử dụng cấu trúc mới coa_account_rule_types
// func (s *OptionsSerive) GetRuleTypes(ctx context.Context) ([]*model.AccountRuleOption, error) {
// 	return s.ruleOptionRepo.GetManyByFields(ctx, map[string]any{
// 		"parent_option_id": nil,
// 		"status":           "ACTIVE",
// 	})
// }

// GetRuleGroups - DEPRECATED: Sử dụng cấu trúc mới coa_account_rule_groups
// func (s *OptionsSerive) GetRuleGroups(ctx context.Context, typeId int64) ([]*model.AccountRuleOption, error) {
// 	return s.ruleOptionRepo.GetManyByFields(ctx, map[string]any{
// 		"parent_option_id": typeId,
// 		"status":           "ACTIVE",
// 	})
// }
