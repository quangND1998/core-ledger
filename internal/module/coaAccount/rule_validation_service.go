package coaaccount

import (
	"context"
	"fmt"
	"strings"

	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type RuleValidationService struct {
	logger        logger.CustomLogger
	db            *gorm.DB
	ruleValueRepo repo.RuleValueRepo
}

func NewRuleValidationService(db *gorm.DB, ruleValueRepo repo.RuleValueRepo) *RuleValidationService {
	return &RuleValidationService{
		logger:        logger.NewSystemLog("RuleValidationService"),
		db:            db,
		ruleValueRepo: ruleValueRepo,
	}
}

// ValidateAndBuildCode validates rule inputs and builds code (not account_no)
// account_no is user input, code is generated from rules
func (s *RuleValidationService) ValidateAndBuildCode(ctx context.Context, ruleInput *dto.CoaAccountRuleInput) (string, error) {
	// 1. Load type và groups/steps từ DB
	var ruleType model.CoaAccountRuleType
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
		Where("code = ?", ruleInput.TypeCode).
		First(&ruleType).Error; err != nil {
		return "", fmt.Errorf("rule type %s not found: %w", ruleInput.TypeCode, err)
	}

	// 2. Xác định group (nếu có)
	var selectedGroup *model.CoaAccountRuleGroup
	var steps []model.CoaAccountRuleStep

	if ruleInput.GroupCode != nil && *ruleInput.GroupCode != "" {
		// Có group - tìm group trong type
		for i := range ruleType.Groups {
			if ruleType.Groups[i].Code == *ruleInput.GroupCode {
				selectedGroup = &ruleType.Groups[i]
				steps = selectedGroup.Steps
				break
			}
		}
		if selectedGroup == nil {
			return "", fmt.Errorf("group %s not found for type %s", *ruleInput.GroupCode, ruleInput.TypeCode)
		}
	} else {
		// Không có group - dùng steps trực tiếp từ type (REV, EXP)
		steps = ruleType.Steps
		if len(steps) == 0 {
			return "", fmt.Errorf("no steps found for type %s", ruleInput.TypeCode)
		}
	}

	// 3. Validate số lượng steps
	if len(ruleInput.Steps) != len(steps) {
		return "", fmt.Errorf("expected %d steps, got %d", len(steps), len(ruleInput.Steps))
	}

	// 4. Build account_no từ các giá trị đã chọn
	var parts []string
	
	// Thêm type
	parts = append(parts, ruleType.Code)
	
	// Thêm group separator và group code nếu có
	if selectedGroup != nil {
		parts = append(parts, ruleType.Separator, selectedGroup.Code)
	}

	// 5. Validate và build từng step
	for i, step := range steps {
		if i >= len(ruleInput.Steps) {
			return "", fmt.Errorf("missing step at order %d", step.StepOrder)
		}

		stepInput := ruleInput.Steps[i]
		
		// Validate step_id và step_order
		if stepInput.StepID != step.ID || stepInput.StepOrder != step.StepOrder {
			return "", fmt.Errorf("step mismatch at order %d", step.StepOrder)
		}

		var value string

		if step.Type == "SELECT" {
			// Validate SELECT step
			if stepInput.ValueID == nil || stepInput.Value == nil {
				return "", fmt.Errorf("step %d (SELECT) requires value_id and value", step.StepOrder)
			}

			// Validate value exists trong rule_values
			ruleValue, err := s.ruleValueRepo.GetByID(ctx, int64(*stepInput.ValueID))
			if err != nil {
				return "", fmt.Errorf("invalid value_id %d for step %d: %w", *stepInput.ValueID, step.StepOrder, err)
			}

			// Validate value matches
			if ruleValue.Value != *stepInput.Value {
				return "", fmt.Errorf("value mismatch for step %d", step.StepOrder)
			}

			// Validate category
			if step.CategoryID != nil && ruleValue.CategoryID != uint(*step.CategoryID) {
				return "", fmt.Errorf("value does not belong to category for step %d", step.StepOrder)
			}

			value = ruleValue.Value
		} else if step.Type == "TEXT" {
			// Validate TEXT step
			if stepInput.Value == nil || *stepInput.Value == "" {
				return "", fmt.Errorf("step %d (TEXT) requires value", step.StepOrder)
			}

			if stepInput.InputCode == nil || *stepInput.InputCode != *step.InputCode {
				return "", fmt.Errorf("input_code mismatch for step %d", step.StepOrder)
			}

			value = *stepInput.Value
		}

		// Thêm separator và value
		// Nếu separator rỗng, chỉ thêm value (thường là step cuối cùng)
		if step.Separator != "" {
			parts = append(parts, step.Separator, value)
		} else {
			parts = append(parts, value)
		}
	}

	// 6. Join tất cả parts thành code
	code := strings.Join(parts, "")

	return code, nil
}

// ParseCodeToRuleInput parses code to extract rule inputs (for edit)
// code format: TYPE[:GROUP][:STEP1.STEP2...]
// Ví dụ: ASSET:FLOAT:VND.Cobo.DETAILS hoặc REV:KINDS_OF_REVENUE.CURRENCY
func (s *RuleValidationService) ParseCodeToRuleInput(ctx context.Context, code string) (*dto.CoaAccountRuleInput, error) {
	// Tách theo separator ":" để lấy type và phần còn lại
	colonIndex := strings.Index(code, ":")
	if colonIndex == -1 {
		return nil, fmt.Errorf("invalid code format: missing type separator")
	}

	typeCode := code[:colonIndex]
	remaining := code[colonIndex+1:]
	
	// Load type từ DB
	var ruleType model.CoaAccountRuleType
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
		Where("code = ?", typeCode).
		First(&ruleType).Error; err != nil {
		return nil, fmt.Errorf("rule type %s not found: %w", typeCode, err)
	}

	ruleInput := &dto.CoaAccountRuleInput{
		TypeCode: typeCode,
		Steps:    []dto.CoaAccountRuleStepInput{},
	}

	// Xác định có group hay không
	var selectedGroup *model.CoaAccountRuleGroup
	var steps []model.CoaAccountRuleStep
	var stepValues []string

	// Kiểm tra xem có group không (tìm separator ":" tiếp theo)
	nextColonIndex := strings.Index(remaining, ":")
	if nextColonIndex != -1 {
		// Có group: TYPE:GROUP:STEP1.STEP2...
		groupCode := remaining[:nextColonIndex]
		stepPart := remaining[nextColonIndex+1:]
		
		for i := range ruleType.Groups {
			if ruleType.Groups[i].Code == groupCode {
				selectedGroup = &ruleType.Groups[i]
				steps = selectedGroup.Steps
				ruleInput.GroupCode = &groupCode
				break
			}
		}
		
		if selectedGroup == nil {
			return nil, fmt.Errorf("group %s not found for type %s", groupCode, typeCode)
		}
		
		// Parse step values theo separator của từng step
		stepValues = s.parseStepValues(stepPart, steps)
	} else {
		// Không có group (REV, EXP): TYPE:STEP1.STEP2...
		steps = ruleType.Steps
		stepValues = s.parseStepValues(remaining, steps)
	}

	// Parse từng step value
	for i, step := range steps {
		if i >= len(stepValues) {
			return nil, fmt.Errorf("missing value for step %d", step.StepOrder)
		}

		value := stepValues[i]
		stepInput := dto.CoaAccountRuleStepInput{
			StepID:    step.ID,
			StepOrder: step.StepOrder,
		}

		if step.Type == "SELECT" {
			// Tìm value trong rule_values
			var ruleValues []model.RuleValue
			if step.CategoryID != nil {
				if err := s.db.WithContext(ctx).
					Where("category_id = ? AND value = ? AND is_delete = false", *step.CategoryID, value).
					Find(&ruleValues).Error; err != nil {
					return nil, fmt.Errorf("failed to find rule value: %w", err)
				}
				
				if len(ruleValues) == 0 {
					return nil, fmt.Errorf("value %s not found for step %d", value, step.StepOrder)
				}
				
				valueID := uint64(ruleValues[0].ID)
				stepInput.ValueID = &valueID
				stepInput.Value = &value
				if step.CategoryCode != nil {
					stepInput.CategoryCode = step.CategoryCode
				}
			}
		} else if step.Type == "TEXT" {
			// TEXT input
			stepInput.Value = &value
			if step.InputCode != nil {
				stepInput.InputCode = step.InputCode
			}
		}

		ruleInput.Steps = append(ruleInput.Steps, stepInput)
	}

	return ruleInput, nil
}

// parseStepValues parses step values từ string dựa trên separators của từng step
// Format: STEP1[SEP1]STEP2[SEP2]STEP3...
func (s *RuleValidationService) parseStepValues(stepPart string, steps []model.CoaAccountRuleStep) []string {
	if len(steps) == 0 {
		return []string{}
	}

	var values []string
	current := stepPart

	for i, step := range steps {
		if i == len(steps)-1 {
			// Step cuối cùng: lấy toàn bộ phần còn lại (bất kể separator)
			if current != "" {
				values = append(values, current)
			} else {
				values = append(values, "")
			}
			break
		}

		// Tìm separator của step hiện tại (dùng để tách step tiếp theo)
		separator := step.Separator
		if separator == "" {
			// Nếu separator rỗng, không thể tách được - lấy toàn bộ và break
			if current != "" {
				values = append(values, current)
			} else {
				values = append(values, "")
			}
			// Fill các step còn lại với empty string
			for j := i + 1; j < len(steps); j++ {
				values = append(values, "")
			}
			break
		}

		// Tìm vị trí separator
		sepIndex := strings.Index(current, separator)
		if sepIndex == -1 {
			// Không tìm thấy separator, lấy phần còn lại và break
			if current != "" {
				values = append(values, current)
			} else {
				values = append(values, "")
			}
			// Fill các step còn lại với empty string
			for j := i + 1; j < len(steps); j++ {
				values = append(values, "")
			}
			break
		}

		// Lấy value trước separator
		value := current[:sepIndex]
		values = append(values, value)

		// Cập nhật current để tiếp tục parse (bỏ qua separator)
		current = current[sepIndex+len(separator):]
	}

	return values
}

