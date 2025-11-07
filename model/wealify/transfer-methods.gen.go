package model

import (
	"time"
)

const TableNameTransferMethod = "transfer-methods"

// TransferMethod mapped from table <transfer-methods>
type TransferMethod struct {
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status      bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	FileID      string    `gorm:"column:file_id" json:"file_id"`
}

// TableName TransferMethod's table name
func (*TransferMethod) TableName() string {
	return TableNameTransferMethod
}
