package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Entry struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	JournalID   uint64          `gorm:"not null;index:idx_entries_journal_id" json:"journal_id"`
	LineNo      int             `gorm:"not null" json:"line_no"`
	AccountID   uint64          `gorm:"not null;index:idx_entries_account_id" json:"account_id"`
	DC          string          `gorm:"type:char(1);not null;check:dc IN ('D','C');index:idx_entries_dc" json:"dc"`
	Amount      decimal.Decimal `gorm:"type:numeric(28,8);not null;check:amount>=0" json:"amount"`
	AmountAtoms *int64          `json:"amount_atoms,omitempty"`
	Meta        map[string]any  `gorm:"type:jsonb" json:"meta,omitempty"`
	Memo        *string         `gorm:"type:varchar(256)" json:"memo,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	TenantID   *string `gorm:"type:varchar(36);index:idx_entries_tenant_id" json:"tenant_id,omitempty"`
	LedgerCode *string `gorm:"type:varchar(32);index:idx_entries_ledger_code" json:"ledger_code,omitempty"`
	BatchID    *string `gorm:"type:varchar(36);index:idx_entries_batch_id" json:"batch_id,omitempty"`

	// Quan hệ
	Journal *Journal    `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	Account *CoaAccount `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	
}

// TableName đặt tên bảng rõ ràng
func (Entry) TableName() string {
	return "entries"
}

// BeforeCreate hook tự động set CreatedAt và UpdatedAt
func (e *Entry) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	return
}

// BeforeUpdate hook tự động set UpdatedAt
func (e *Entry) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()
	return
}
