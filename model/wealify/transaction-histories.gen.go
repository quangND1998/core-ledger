package model

import (
	"core-ledger/model/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameTransactionHistorie = "transaction-histories"

// TransactionHistorie mapped from table <transaction-histories>
type TransactionHistorie struct {
	CreatedAt         time.Time                     `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt         time.Time                     `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status            int32                         `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted         int32                         `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID                string                        `gorm:"column:id;primaryKey" json:"id"`
	Type              TransactionHistoryType        `gorm:"column:type;not null;default:CREATED" json:"type"`
	TransactionStatus enum.TransactionStatus        `gorm:"column:transaction_status" json:"transaction_status"`
	Amount            float64                       `gorm:"column:amount" json:"amount"`
	TransactionID     string                        `gorm:"column:transaction_id" json:"transaction_id"`
	EmployeeID        *int64                        `gorm:"column:employee_id" json:"employee_id"`
	CustomerID        *int64                        `gorm:"column:customer_id" json:"customer_id"`
	PostBalance       float64                       `gorm:"column:post_balance;not null" json:"post_balance"`
	PreBalance        float64                       `gorm:"column:pre_balance;not null" json:"pre_balance"`
	Source            *TransactionHistorySourceType `gorm:"column:source" json:"source"`
}

// TableName TransactionHistorie's table name
func (*TransactionHistorie) TableName() string {
	return TableNameTransactionHistorie
}

func (t *TransactionHistorie) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *TransactionHistorie) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}
