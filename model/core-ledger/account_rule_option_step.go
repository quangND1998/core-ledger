package model

import "time"

type AccountRuleOptionStep struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OptionID   uint64     `gorm:"not null" json:"option_id"`
	StepOrder  int        `gorm:"not null" json:"step_order"`
	CategoryID *uint64    `json:"category_id,omitempty"`
	InputCode  *string    `gorm:"type:varchar(64)" json:"input_code,omitempty"`
	InputLabel *string    `gorm:"type:varchar(128)" json:"input_label,omitempty"`
	InputType  string     `gorm:"type:varchar(16);default:'SELECT'" json:"input_type"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AccountRuleOptionStep) TableName() string {
	return "account_rule_option_steps"
}


