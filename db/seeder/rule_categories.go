package seeder

import (
	model "core-ledger/model/core-ledger"
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeederRuleCategories(db *gorm.DB) error {
	rule_categories := []model.RuleCategory{
		{
			Code:     "CURRENCY",
			Name:     "Currencies",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
		{
			Code:     "PROVIDER",
			Name:     "Providers",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
		{
			Code:     "BANK_NAME",
			Name:     "Bank Names",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
		{
			Code:     "NETWORK",
			Name:     "Networks",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
		{
			Code:     "KINDS_OF_REVENUE",
			Name:     "Kinds of Revenue",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
		{
			Code:     "KINDS_OF_EXPENSE",
			Name:     "Kinds of Expense",
			Metadata: datatypes.JSONMap{"separator": "."},
		},
	}

	for _, rc := range rule_categories {
		var existing model.RuleCategory
		err := db.Where("code = ?", rc.Code).First(&existing).Error
		if err == nil {
			// Đã tồn tại, merge metadata để cập nhật separator
			mergedMetadata := make(datatypes.JSONMap)
			if existing.Metadata != nil {
				for k, v := range existing.Metadata {
					mergedMetadata[k] = v
				}
			}
			if rc.Metadata != nil {
				for k, v := range rc.Metadata {
					mergedMetadata[k] = v
				}
			}
			if err := db.Model(&existing).Update("metadata", mergedMetadata).Error; err != nil {
				fmt.Println("Error updating rule category metadata:", rc.Code, err)
				continue
			}
			fmt.Println("Updated rule category metadata:", rc.Code)
			continue
		}
		if err := db.Create(&rc).Error; err != nil {
			fmt.Println("Error creating rule category:", rc.Code, err)
			continue
		}
		fmt.Println("Created rule category:", rc.Code)
	}

	return nil
}
