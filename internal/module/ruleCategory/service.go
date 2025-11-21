package ruleCategory

import (
	"core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RuleCateogySerive struct {
	db *gorm.DB

	logger     logger.CustomLogger
	dispatcher queue.Dispatcher
}

func NewRuleCateogySerive(dispatcher queue.Dispatcher, db *gorm.DB) *RuleCateogySerive {
	return &RuleCateogySerive{
		db:         db,
		logger:     logger.NewSystemLog("RuleCateogySerive"),
		dispatcher: dispatcher,
	}
}

// GetCoaAccountRules lấy toàn bộ cấu trúc rules để tạo mã COA account
func (s *RuleCateogySerive) GetCoaAccountRules(c *gin.Context) ([]dto.CoaAccountRuleTypeResp, error) {
	ctx := c.Request.Context()

	// Load tất cả types với groups và steps
	var types []model.CoaAccountRuleType
	if err := s.db.WithContext(ctx).
		Preload("Groups", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("Groups.Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_order ASC")
		}).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Where("group_id IS NULL").Order("step_order ASC")
		}).
		Order("sort_order ASC").
		Find(&types).Error; err != nil {
		return nil, err
	}

	// Load tất cả rule values để map
	var allValues []model.RuleValue
	if err := s.db.WithContext(ctx).
		Where("is_delete = false").
		Order("sort_order ASC").
		Find(&allValues).Error; err != nil {
		return nil, err
	}

	// Map values theo category_id
	valueMap := make(map[uint][]dto.RuleValueResp)
	for _, v := range allValues {
		valueMap[v.CategoryID] = append(valueMap[v.CategoryID], dto.RuleValueResp{
			ID:    uint64(v.ID),
			Name:  v.Name,
			Value: v.Value,
		})
	}

	// Build response
	var result []dto.CoaAccountRuleTypeResp
	for _, typ := range types {
		typeResp := dto.CoaAccountRuleTypeResp{
			ID:        typ.ID,
			Code:      typ.Code,
			Name:      typ.Name,
			Separator: typ.Separator,
			Groups:    []dto.CoaAccountRuleGroupResp{},
		}

		// Process groups
		for _, group := range typ.Groups {
			groupResp := dto.CoaAccountRuleGroupResp{
				ID:        group.ID,
				Code:      group.Code,
				Name:      group.Name,
				InputType: group.InputType,
				Separator: group.Separator,
				Steps:     []dto.CoaAccountRuleStepResp{},
			}

			// Process steps trong group
			for _, step := range group.Steps {
				stepResp := s.buildStepResp(step, valueMap)
				groupResp.Steps = append(groupResp.Steps, stepResp)
			}

			typeResp.Groups = append(typeResp.Groups, groupResp)
		}

		// Nếu type có steps trực tiếp (không có group) - như REV, EXP
		// Tạo group ảo "KIND" để giữ format JSON giống như user cung cấp
		if len(typ.Steps) > 0 {
			groupResp := dto.CoaAccountRuleGroupResp{
				ID:        0, // Không có ID vì không có group trong DB
				Code:      "KIND",
				Name:      "Kind",
				InputType: "SELECT",
				Separator: ":",
				Steps:     []dto.CoaAccountRuleStepResp{},
			}

			for _, step := range typ.Steps {
				stepResp := s.buildStepResp(step, valueMap)
				groupResp.Steps = append(groupResp.Steps, stepResp)
			}

			typeResp.Groups = append(typeResp.Groups, groupResp)
		}

		result = append(result, typeResp)
	}

	return result, nil
}

func (s *RuleCateogySerive) buildStepResp(step model.CoaAccountRuleStep, valueMap map[uint][]dto.RuleValueResp) dto.CoaAccountRuleStepResp {
	stepResp := dto.CoaAccountRuleStepResp{
		StepID:    step.ID,
		StepOrder: step.StepOrder,
		Type:      step.Type,
		Separator: step.Separator,
	}

	if step.Label != nil {
		stepResp.Label = *step.Label
	}

	if step.Type == "SELECT" && step.CategoryID != nil {
		// SELECT type: có category và values
		if step.CategoryCode != nil {
			stepResp.CategoryCode = *step.CategoryCode
		}
		stepResp.InputType = "SELECT"
		if step.CategoryID != nil {
			catID := uint(*step.CategoryID)
			stepResp.Values = valueMap[catID]
		}
	} else if step.Type == "TEXT" && step.InputCode != nil {
		// TEXT type: có input_code
		stepResp.InputCode = *step.InputCode
		stepResp.InputType = step.Type
	}

	return stepResp
}
