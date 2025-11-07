package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const TableNameConfirmTransaction = "confirm-transactions"

// ConfirmTransaction mapped from table <confirm-transactions>
type ConfirmTransaction struct {
	CreatedAt            time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt            time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status               int32           `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted            int32           `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID                   string          `gorm:"column:id;primaryKey" json:"id"`
	ConfirmTransactionID string          `gorm:"column:confirm_transaction_id;not null" json:"confirm_transaction_id"`
	TransactionID        string          `gorm:"column:transaction_id" json:"transaction_id"`
	Transaction          *Transaction    `gorm:"foreignKey:TransactionID" json:"transaction"`
	RawData              *datatypes.JSON `gorm:"column:raw_data" json:"raw_data"`
}

// TableName ConfirmTransaction's table name
func (*ConfirmTransaction) TableName() string {
	return TableNameConfirmTransaction
}

func (ct *ConfirmTransaction) BeforeSave(tx *gorm.DB) (err error) {
	if ct.ID == "" {
		ct.ID = uuid.NewString()
	}
	return
}

func (ct *ConfirmTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	if ct.ID == "" {
		ct.ID = uuid.NewString()
	}
	return
}
