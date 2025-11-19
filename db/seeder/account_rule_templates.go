package seeder

import (
	"errors"
	"fmt"
	"strings"

	model "core-ledger/model/core-ledger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeederAccountRuleTemplates(db *gorm.DB) error {
	layerSeeds := []model.AccountRuleLayer{
		{Code: "TYPE", Name: "Layer 1 (Type)", LayerIndex: 1, Status: "ACTIVE"},
		{Code: "GROUP", Name: "Layer 2 (Group)", LayerIndex: 2, Status: "ACTIVE"},
		{Code: "CURRENCY", Name: "Layer 3", LayerIndex: 3, Status: "ACTIVE"},
		{Code: "DETAILS", Name: "Layer 4", LayerIndex: 4, Status: "ACTIVE", Metadata: datatypes.JSONMap{
			"input_type":  "text",
			"placeholder": "Nhập DETAILS",
		}},
		{Code: "NETWORK", Name: "Layer 5", LayerIndex: 5, Status: "ACTIVE"},
	}

	layerMap := make(map[string]*model.AccountRuleLayer)
	for _, seed := range layerSeeds {
		layer, err := upsertAccountRuleLayer(db, seed)
		if err != nil {
			return fmt.Errorf("seed layer %s: %w", seed.Code, err)
		}
		layerMap[seed.Code] = layer
	}

	optionMap := make(map[string]uint64)

	typeSeed := []model.AccountRuleOption{
		{LayerID: layerMap["TYPE"].ID, Code: "ASSET", Name: "Asset"},
		{LayerID: layerMap["TYPE"].ID, Code: "LIAB", Name: "Liability", InputType: "TEXT"},
		{LayerID: layerMap["TYPE"].ID, Code: "REV", Name: "Revenue"},
		{LayerID: layerMap["TYPE"].ID, Code: "EXP", Name: "Expense"},
	}

	for _, seed := range typeSeed {
		opt, err := upsertAccountRuleOption(db, seed)
		if err != nil {
			return fmt.Errorf("seed type %s: %w", seed.Code, err)
		}
		optionMap["TYPE:"+seed.Code] = opt.ID
	}

	groupSeeds := []struct {
		ParentType string
		Code       string
		Name       string
	}{
		{
			ParentType: "ASSET",
			Code:       "FLOAT",
			Name:       "Float",
		},
		{
			ParentType: "ASSET",
			Code:       "BANK",
			Name:       "Bank",
		},
		{
			ParentType: "ASSET",
			Code:       "CUSTODY",
			Name:       "Custody",
		},
		{
			ParentType: "ASSET",
			Code:       "CLEARING",
			Name:       "Clearing",
		},
		{
			ParentType: "LIAB",
			Code:       "DETAILS",
			Name:       "Details",
		},
		{
			ParentType: "REV",
			Code:       "KIND",
			Name:       "Kind",
		},
		{
			ParentType: "EXP",
			Code:       "KIND",
			Name:       "Kind",
		},
	}

	groupLayerID := layerMap["GROUP"].ID
	for _, seed := range groupSeeds {
		parentID := optionMap["TYPE:"+seed.ParentType]
		optSeed := model.AccountRuleOption{
			LayerID:        groupLayerID,
			ParentOptionID: &parentID,
			Code:           seed.Code,
			Name:           seed.Name,
		}
		opt, err := upsertAccountRuleOption(db, optSeed)
		if err != nil {
			return fmt.Errorf("seed group %s: %w", seed.Code, err)
		}
		key := fmt.Sprintf("GROUP:%s:%s", seed.ParentType, seed.Code)
		optionMap[key] = opt.ID
	}

	categoryMap, err := buildRuleCategoryMap(db)
	if err != nil {
		return err
	}

	stepDefs := map[string][]stepDefinition{
		"GROUP:ASSET:FLOAT": {
			{CategoryCode: "CURRENCY"},
			{CategoryCode: "PROVIDER"},
			{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT"},
		},
		"GROUP:ASSET:BANK": {
			{CategoryCode: "CURRENCY"},
			{CategoryCode: "BANK_NAME"},
			{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT"},
		},
		"GROUP:ASSET:CUSTODY": {
			{CategoryCode: "CURRENCY"},
			{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT"},
			{CategoryCode: "NETWORK"},
		},
		"GROUP:ASSET:CLEARING": {
			{CategoryCode: "CURRENCY"},
		},
		"GROUP:LIAB:DETAILS": {
			{CategoryCode: "CURRENCY"},
		},
		"GROUP:REV:KIND": {
			{CategoryCode: "KINDS_OF_REVENUE"},
			{CategoryCode: "CURRENCY"},
		},
		"GROUP:EXP:KIND": {
			{CategoryCode: "KINDS_OF_EXPENSE"},
			{CategoryCode: "CURRENCY"},
		},
	}

	for optionKey, defs := range stepDefs {
		optionID, ok := optionMap[optionKey]
		if !ok {
			return fmt.Errorf("missing option %s for step seeding", optionKey)
		}
		if err := replaceOptionSteps(db, optionID, defs, categoryMap); err != nil {
			return err
		}
	}

	fmt.Println("Seeded account rule templates")
	return nil
}

func upsertAccountRuleLayer(db *gorm.DB, seed model.AccountRuleLayer) (*model.AccountRuleLayer, error) {
	var existing model.AccountRuleLayer
	err := db.Where("code = ?", seed.Code).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if seed.Status == "" {
			seed.Status = "ACTIVE"
		}
		if seed.Metadata == nil {
			seed.Metadata = datatypes.JSONMap{}
		}
		if err := db.Create(&seed).Error; err != nil {
			return nil, err
		}
		return &seed, nil
	}
	updates := map[string]interface{}{
		"name":        seed.Name,
		"layer_index": seed.LayerIndex,
		"status":      seed.Status,
		"description": seed.Description,
		"metadata":    seed.Metadata,
	}
	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func upsertAccountRuleOption(db *gorm.DB, seed model.AccountRuleOption) (*model.AccountRuleOption, error) {
	query := db.Where("layer_id = ? AND code = ?", seed.LayerID, seed.Code)
	if seed.ParentOptionID == nil {
		query = query.Where("parent_option_id IS NULL")
	} else {
		query = query.Where("parent_option_id = ?", *seed.ParentOptionID)
	}

	var existing model.AccountRuleOption
	err := query.First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if seed.Status == "" {
			seed.Status = "ACTIVE"
		}
		if seed.Metadata == nil {
			seed.Metadata = datatypes.JSONMap{}
		}
		if err := db.Create(&seed).Error; err != nil {
			return nil, err
		}
		return &seed, nil
	}

	updates := map[string]interface{}{
		"name":       seed.Name,
		"status":     defaultStatus(seed.Status, existing.Status),
		"sort_order": seed.SortOrder,
		"metadata":   seed.Metadata,
	}
	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func defaultStatus(candidate, fallback string) string {
	if candidate != "" {
		return candidate
	}
	return fallback
}

type stepDefinition struct {
	CategoryCode string
	InputCode    string
	InputLabel   string
	InputType    string
}

func replaceOptionSteps(db *gorm.DB, optionID uint64, defs []stepDefinition, categoryMap map[string]uint64) error {
	if err := db.Where("option_id = ?", optionID).Delete(&model.AccountRuleOptionStep{}).Error; err != nil {
		return err
	}
	for idx, def := range defs {
		step := model.AccountRuleOptionStep{
			OptionID:  optionID,
			StepOrder: idx + 1,
		}
		if def.CategoryCode != "" {
			catID, ok := categoryMap[def.CategoryCode]
			if !ok {
				return fmt.Errorf("rule category %s not found", def.CategoryCode)
			}
			step.CategoryID = uint64Pointer(catID)
			step.InputType = "SELECT"
		} else {
			if def.InputCode == "" {
				return fmt.Errorf("option %d step %d missing category or input", optionID, idx+1)
			}
			step.InputCode = stringPointer(strings.ToUpper(def.InputCode))
			if def.InputLabel != "" {
				step.InputLabel = stringPointer(def.InputLabel)
			}
			if def.InputType != "" {
				step.InputType = strings.ToUpper(def.InputType)
			} else {
				step.InputType = "TEXT"
			}
		}
		if err := db.Create(&step).Error; err != nil {
			return err
		}
	}
	return nil
}

func buildRuleCategoryMap(db *gorm.DB) (map[string]uint64, error) {
	var categories []model.RuleCategory
	if err := db.Select("id", "code").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("load rule categories: %w", err)
	}
	result := make(map[string]uint64, len(categories))
	for _, cat := range categories {
		result[strings.ToUpper(cat.Code)] = uint64(cat.ID)
	}
	return result, nil
}

func stringPointer(s string) *string {
	if s == "" {
		return nil
	}
	val := s
	return &val
}

func uint64Pointer(v uint64) *uint64 {
	val := v
	return &val
}
