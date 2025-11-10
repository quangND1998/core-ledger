package model

import (
	"time"

	"gorm.io/gorm"
)

type TransactionLog struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	AggregateType string         `gorm:"type:varchar(32);not null;index:idx_transaction_logs_aggregate" json:"aggregate_type"`
	AggregateID   uint64         `gorm:"not null;index:idx_transaction_logs_aggregate" json:"aggregate_id"`
	EventType     string         `gorm:"type:varchar(64);not null" json:"event_type"`
	EventKey      string         `gorm:"type:varchar(191);not null" json:"event_key"`
	PartitionKey  string         `gorm:"type:varchar(191);not null;index:idx_transaction_logs_partition_key" json:"partition_key"`
	Payload       map[string]any `gorm:"type:jsonb;not null" json:"payload"`
	Headers       map[string]any `gorm:"type:jsonb" json:"headers,omitempty"`
	Status        string         `gorm:"type:varchar(16);not null;default:'PENDING';check:status IN ('PENDING','PUBLISHED','FAILED','DEAD');index:idx_transaction_logs_status" json:"status"`
	Attempts      int            `gorm:"not null;default:0" json:"attempts"`
	NextAttemptAt *time.Time     `gorm:"index:idx_transaction_logs_next_attempt" json:"next_attempt_at,omitempty"`
	LastAttemptAt *time.Time     `json:"last_attempt_at,omitempty"`
	PublishedAt   *time.Time     `json:"published_at,omitempty"`
	ErrorLast     *string        `gorm:"type:text" json:"error_last,omitempty"`
	TenantID      string         `gorm:"type:varchar(36);not null;index:idx_transaction_logs_tenant" json:"tenant_id"`
	LedgerCode    *string        `gorm:"type:varchar(32);index:idx_transaction_logs_ledger" json:"ledger_code,omitempty"`
	Seq           *uint64        `gorm:"index:idx_transaction_logs_seq" json:"seq,omitempty"`
	CreatedAt     time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
}

func (TransactionLog) TableName() string {
	return "transaction_logs"
}

// Hooks
func (t *TransactionLog) BeforeCreate(tx *gorm.DB) (err error) {
	t.CreatedAt = time.Now()
	return
}

func (t *TransactionLog) BeforeUpdate(tx *gorm.DB) (err error) {
	// Nếu cần, update các trường tracking
	return
}
