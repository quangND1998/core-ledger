package model

import (
	"time"
)

const TableNameCallingCode = "calling-codes"

// CallingCode mapped from table <calling-codes>
type CallingCode struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status    bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Code      int64     `gorm:"column:code;not null" json:"code"`
}

// TableName CallingCode's table name
func (*CallingCode) TableName() string {
	return TableNameCallingCode
}
