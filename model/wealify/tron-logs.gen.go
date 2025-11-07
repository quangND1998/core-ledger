package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameTronLog = "tron-logs"

// TronLog mapped from table <tron-logs>
type TronLog struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Token     string    `gorm:"column:token;not null;default:USDT" json:"token"`
	To        string    `gorm:"column:to;not null" json:"to"`
	From      string    `gorm:"column:from;not null" json:"from"`
	Amount    float64   `gorm:"column:amount;not null" json:"amount"`
	TxHash    string    `gorm:"column:txHash;not null" json:"txHash"`
	Timestamp int64     `gorm:"column:timestamp;not null" json:"timestamp"`
	Mapped    bool      `gorm:"column:mapped;not null" json:"mapped"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status    bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	Chain     string    `gorm:"column:chain;not null;default:TRX" json:"chain"`
}

// TableName TronLog's table name
func (*TronLog) TableName() string {
	return TableNameTronLog
}

func (t *TronLog) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *TronLog) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}
