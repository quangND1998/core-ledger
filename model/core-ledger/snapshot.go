package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Snapshot struct {
	ID             uint64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AsOfDate       time.Time       `gorm:"type:date;not null;index:idx_snapshots_as_of_date" json:"as_of_date"`
	AccountID      uint64          `gorm:"not null;index:idx_snapshots_account_id" json:"account_id"`
	AccountCode    string          `gorm:"type:varchar(128);not null;index:idx_snapshots_account_code" json:"account_code"`
	Currency       string          `gorm:"type:char(8);not null" json:"currency"`
	OpeningBalance decimal.Decimal `gorm:"type:numeric(28,8);not null" json:"opening_balance"`
	DebitTotal     decimal.Decimal `gorm:"type:numeric(28,8);not null" json:"debit_total"`
	CreditTotal    decimal.Decimal `gorm:"type:numeric(28,8);not null" json:"credit_total"`
	Movement       decimal.Decimal `gorm:"type:numeric(28,8);not null" json:"movement"`
	ClosingBalance decimal.Decimal `gorm:"type:numeric(28,8);not null" json:"closing_balance"`
	EntryCount     int             `gorm:"not null" json:"entry_count"`
	LedgerCode     *string         `gorm:"type:varchar(32)" json:"ledger_code,omitempty"`
	TenantID       *string         `gorm:"type:varchar(36);index:idx_snapshots_tenant_id" json:"tenant_id,omitempty"`
	Status         string          `gorm:"type:varchar(16);default:'DRAFT';check:status IN ('DRAFT','LOCKED')" json:"status"`
	Hash           string          `gorm:"type:char(64);not null" json:"hash"`
	Meta           map[string]any  `gorm:"type:jsonb" json:"meta,omitempty"`
	CreatedAt      time.Time       `gorm:"column:created_at" json:"created_at"`
	CreatedBy      *string         `gorm:"type:varchar(64)" json:"created_by,omitempty"`

	// Quan hệ
	Account *CoaAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

func (Snapshot) TableName() string {
	return "snapshots"
}

// Hooks
func (s *Snapshot) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreatedAt = time.Now()
	return
}

// Nếu sau này muốn track update
func (s *Snapshot) BeforeUpdate(tx *gorm.DB) (err error) {
	// hiện tại bảng chưa có updated_at, nhưng có thể thêm nếu cần
	return
}
