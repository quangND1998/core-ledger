package model

import (
	"time"
)

const TableNameRespondent = "respondents"

// Respondent mapped from table <respondents>
type Respondent struct {
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status     int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted  int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	Data       string    `gorm:"column:data;not null" json:"data"`
	RequestID  string    `gorm:"column:request_id" json:"request_id"`
	CustomerID int32     `gorm:"column:customer_id" json:"customer_id"`
}

// TableName Respondent's table name
func (*Respondent) TableName() string {
	return TableNameRespondent
}
