package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameCountrie = "countries"

// Countrie mapped from table <countries>
type Countrie struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status    bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Code      int64     `gorm:"column:code" json:"code"`
	IsoCode2  string    `gorm:"column:iso_code2;not null" json:"iso_code2"`
	IsoCode3  string    `gorm:"column:iso_code3;not null" json:"iso_code3"`
	FileID    *string   `gorm:"column:file_id" json:"file_id"`
}

// TableName Countrie's table name
func (*Countrie) TableName() string {
	return TableNameCountrie
}

func (c *Countrie) BeforeSave(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}

func (c *Countrie) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Countrie) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
