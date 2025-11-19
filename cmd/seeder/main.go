package main

import (
	"core-ledger/db/seeder"
	"core-ledger/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

func main() {
	db := database.Instance()

	seeders := []func(*gorm.DB) error{
		seeder.SeederRuleCategories,
		seeder.SeederAccountRuleTemplates,
	}

	for _, s := range seeders {
		if err := s(db); err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Seeder completed!")
}
