package model

import (
	"time"
)

const TableNameWorldFirstActivitie = "world-first-activities"

// WorldFirstActivitie mapped from table <world-first-activities>
type WorldFirstActivitie struct {
	ActivityID     string    `gorm:"column:activity_id;primaryKey" json:"activity_id"`
	Amount         float32   `gorm:"column:amount;not null" json:"amount"`
	CurrencySymbol string    `gorm:"column:currency_symbol" json:"currency_symbol"`
	Detail         string    `gorm:"column:detail" json:"detail"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
}

// TableName WorldFirstActivitie's table name
func (*WorldFirstActivitie) TableName() string {
	return TableNameWorldFirstActivitie
}
