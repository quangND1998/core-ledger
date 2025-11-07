package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameLog = "logs"

// Log mapped from table <logs>
type Log struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Log       string    `gorm:"column:log" json:"log"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
}

// TableName Log's table name
func (*Log) TableName() string {
	return TableNameLog
}

func (l *Log) BeforeSave(tx *gorm.DB) (err error) {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	return
}

func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	return
}
