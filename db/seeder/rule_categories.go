package seeder

import (
	model "core-ledger/model/core-ledger"
	"fmt"

	"gorm.io/gorm"
)

func SeederRuleCategories(db *gorm.DB) error {
	rule_categories := []model.RuleCategory{
		{Code: "CURRENCY", Name: "Currencies"},
		{Code: "PROVIDER", Name: "Providers"},
		{Code: "BANK_NAME", Name: "Bank Names"},
		{Code: "NETWORK", Name: "Networks"},
		{Code: "KINDS_OF_REVENUE", Name: "Kinds of Revenue"},
		{Code: "KINDS_OF_EXPENSE", Name: "Kinds of Expense"},
	}

	for _, rc := range rule_categories {
		var existing model.RuleCategory
		err := db.Where("code = ?", rc.Code).First(&existing).Error
		if err == nil {
			// đã tồn tại, bỏ qua
			fmt.Println("Rule category already exists:", rc.Code)
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
