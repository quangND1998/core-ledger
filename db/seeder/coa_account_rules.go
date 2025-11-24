package seeder

import (
	model "core-ledger/model/core-ledger"
	"errors"
	"fmt"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SeederCoaAccountRules seed data cho cấu trúc rule mới
func SeederCoaAccountRules(db *gorm.DB) error {
	// 1. Seed TYPES
	typeSeeds := []model.CoaAccountRuleType{
		{Code: "ASSET", Name: "Asset", Separator: ":", SortOrder: 1},
		{Code: "LIAB", Name: "Liability", Separator: ":", SortOrder: 2},
		{Code: "REV", Name: "Revenue", Separator: ":", SortOrder: 3},
		{Code: "EXP", Name: "Expense", Separator: ":", SortOrder: 4},
	}

	typeMap := make(map[string]uint64)
	for _, seed := range typeSeeds {
		typ, err := upsertCoaAccountRuleType(db, seed)
		if err != nil {
			return fmt.Errorf("seed type %s: %w", seed.Code, err)
		}
		typeMap[seed.Code] = typ.ID
	}

	// 2. Seed GROUPS
	groupSeeds := []struct {
		TypeCode  string
		Code      string
		Name      string
		Separator string
		InputType string
		SortOrder int
	}{
		// ASSET groups
		{"ASSET", "FLOAT", "Float", ":", "SELECT", 1},
		{"ASSET", "BANK", "Bank", ":", "SELECT", 2},
		{"ASSET", "CUSTODY", "Custody", ":", "SELECT", 3},
		{"ASSET", "CLEARING", "Clearing", ":", "SELECT", 4},
		// LIAB groups
		{"LIAB", "DETAILS", "Details", ":", "TEXT", 1},
		// REV và EXP groups
		{"REV", "KIND", "Kind", ":", "SELECT", 1},
		{"EXP", "KIND", "Kind", ":", "SELECT", 1},
	}

	groupMap := make(map[string]uint64)
	for _, seed := range groupSeeds {
		typeID := typeMap[seed.TypeCode]
		group := model.CoaAccountRuleGroup{
			TypeID:    typeID,
			Code:      seed.Code,
			Name:      seed.Name,
			Separator: seed.Separator,
			InputType: seed.InputType,
			SortOrder: seed.SortOrder,
		}
		g, err := upsertCoaAccountRuleGroup(db, group)
		if err != nil {
			return fmt.Errorf("seed group %s:%s: %w", seed.TypeCode, seed.Code, err)
		}
		key := fmt.Sprintf("%s:%s", seed.TypeCode, seed.Code)
		groupMap[key] = g.ID
	}

	// 3. Load rule categories để map
	categoryMap, err := buildRuleCategoryMap(db)
	if err != nil {
		return err
	}

	// 4. Seed STEPS
	// ASSET:FLOAT steps
	if err := seedStepsForGroup(db, typeMap["ASSET"], groupMap["ASSET:FLOAT"], []stepSeed{
		{CategoryCode: "CURRENCY", Separator: ".", StepOrder: 1},
		{CategoryCode: "PROVIDER", Separator: ".", StepOrder: 2},
		{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT", Separator: "", StepOrder: 3},
	}, categoryMap); err != nil {
		return err
	}

	// ASSET:BANK steps
	if err := seedStepsForGroup(db, typeMap["ASSET"], groupMap["ASSET:BANK"], []stepSeed{
		{CategoryCode: "CURRENCY", Separator: ".", StepOrder: 1},
		{CategoryCode: "BANK_NAME", Separator: ".", StepOrder: 2},
		{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT", Separator: "", StepOrder: 3},
	}, categoryMap); err != nil {
		return err
	}

	// ASSET:CUSTODY steps
	if err := seedStepsForGroup(db, typeMap["ASSET"], groupMap["ASSET:CUSTODY"], []stepSeed{
		{CategoryCode: "CURRENCY", Separator: ".", StepOrder: 1},
		{InputCode: "DETAILS", InputLabel: "Details", InputType: "TEXT", Separator: ".", StepOrder: 2},
		{CategoryCode: "NETWORK", Separator: "", StepOrder: 3},
	}, categoryMap); err != nil {
		return err
	}

	// ASSET:CLEARING steps
	if err := seedStepsForGroup(db, typeMap["ASSET"], groupMap["ASSET:CLEARING"], []stepSeed{
		{CategoryCode: "CURRENCY", Separator: "", StepOrder: 1},
	}, categoryMap); err != nil {
		return err
	}

	// LIAB:DETAILS steps
	if err := seedStepsForGroup(db, typeMap["LIAB"], groupMap["LIAB:DETAILS"], []stepSeed{
		{CategoryCode: "CURRENCY", Separator: "", StepOrder: 1},
	}, categoryMap); err != nil {
		return err
	}

	// REV:KIND steps
	if err := seedStepsForGroup(db, typeMap["REV"], groupMap["REV:KIND"], []stepSeed{
		{CategoryCode: "KINDS_OF_REVENUE", Separator: ".", StepOrder: 1},
		{CategoryCode: "CURRENCY", Separator: "", StepOrder: 2},
	}, categoryMap); err != nil {
		return err
	}

	// EXP:KIND steps
	if err := seedStepsForGroup(db, typeMap["EXP"], groupMap["EXP:KIND"], []stepSeed{
		{CategoryCode: "KINDS_OF_EXPENSE", Separator: ".", StepOrder: 1},
		{CategoryCode: "CURRENCY", Separator: "", StepOrder: 2},
	}, categoryMap); err != nil {
		return err
	}

	fmt.Println("Seeded COA account rules successfully")
	return nil
}

type stepSeed struct {
	CategoryCode string
	InputCode    string
	InputLabel   string
	InputType    string
	Separator    string
	StepOrder    int
}

func seedStepsForGroup(db *gorm.DB, typeID uint64, groupID uint64, steps []stepSeed, categoryMap map[string]uint64) error {
	return seedStepsForType(db, typeID, &groupID, steps, categoryMap)
}

func seedStepsForType(db *gorm.DB, typeID uint64, groupID *uint64, steps []stepSeed, categoryMap map[string]uint64) error {
	// Xóa steps cũ
	query := db.Where("type_id = ?", typeID)
	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	} else {
		query = query.Where("group_id IS NULL")
	}
	if err := query.Delete(&model.CoaAccountRuleStep{}).Error; err != nil {
		return err
	}

	// Tạo steps mới
	for _, seed := range steps {
		step := model.CoaAccountRuleStep{
			TypeID:    typeID,
			GroupID:   groupID,
			StepOrder: seed.StepOrder,
			Separator: seed.Separator,
		}

		if seed.CategoryCode != "" {
			// SELECT type
			catID, ok := categoryMap[seed.CategoryCode]
			if !ok {
				return fmt.Errorf("category %s not found", seed.CategoryCode)
			}
			step.Type = "SELECT"
			step.CategoryID = &catID
			code := seed.CategoryCode
			step.CategoryCode = &code
		} else if seed.InputCode != "" {
			// TEXT type
			step.Type = seed.InputType
			if step.Type == "" {
				step.Type = "TEXT"
			}
			step.InputCode = &seed.InputCode
			if seed.InputLabel != "" {
				step.Label = &seed.InputLabel
			}
		}

		if err := db.Create(&step).Error; err != nil {
			return fmt.Errorf("create step: %w", err)
		}
	}

	return nil
}

func upsertCoaAccountRuleType(db *gorm.DB, seed model.CoaAccountRuleType) (*model.CoaAccountRuleType, error) {
	var existing model.CoaAccountRuleType
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

	// Update
	updates := map[string]interface{}{
		"name":       seed.Name,
		"separator":  seed.Separator,
		"status":     seed.Status,
		"sort_order": seed.SortOrder,
	}
	if seed.Metadata != nil {
		updates["metadata"] = seed.Metadata
	}
	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func upsertCoaAccountRuleGroup(db *gorm.DB, seed model.CoaAccountRuleGroup) (*model.CoaAccountRuleGroup, error) {
	var existing model.CoaAccountRuleGroup
	err := db.Where("type_id = ? AND code = ?", seed.TypeID, seed.Code).First(&existing).Error
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

	// Update
	updates := map[string]interface{}{
		"name":       seed.Name,
		"separator":  seed.Separator,
		"input_type": seed.InputType,
		"status":     seed.Status,
		"sort_order": seed.SortOrder,
	}
	if seed.Metadata != nil {
		updates["metadata"] = seed.Metadata
	}
	if err := db.Model(&existing).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &existing, nil
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
