package model

import (
	"time"
)

type Journal struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Ts             time.Time      `gorm:"column:ts;not null;default:now()" json:"ts"`
	Status         string         `gorm:"type:varchar(16);not null;check:status IN ('DRAFT','POSTED','REVERSED')" json:"status"`
	IdempotencyKey string         `gorm:"type:varchar(191);unique;not null" json:"idempotency_key"`
	Currency       string         `gorm:"type:char(8);not null" json:"currency"`
	Source         string         `gorm:"type:varchar(64);not null" json:"source"`
	Memo           *string        `gorm:"type:varchar(256)" json:"memo,omitempty"`
	Meta           map[string]any `gorm:"type:jsonb" json:"meta,omitempty"`
	ReversalOfID   *uint64        `gorm:"column:reversal_of" json:"reversal_of,omitempty"`
	ReversalOf     *Journal       `gorm:"foreignKey:ReversalOfID" json:"reversal_of_journal,omitempty"`
	PostedBy       *string        `gorm:"type:varchar(64)" json:"posted_by,omitempty"`
	PostedAt       *time.Time     `gorm:"column:posted_at" json:"posted_at,omitempty"`
	LockVersion    int            `gorm:"default:0;not null" json:"lock_version"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime;not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime;not null" json:"updated_at"`
	TenantID       *string        `gorm:"type:varchar(36)" json:"tenant_id,omitempty"`
	LedgerCode     *string        `gorm:"type:varchar(32)" json:"ledger_code,omitempty"`
	BatchID        *string        `gorm:"type:varchar(36)" json:"batch_id,omitempty"`
}

// TableName đặt tên bảng rõ ràng
func (Journal) TableName() string {
	return "journals"
}
