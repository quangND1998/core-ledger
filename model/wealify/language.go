package model

import (
	"time"
)

const TableNameLanguage = "languages"

// Language mapped from table <languages>
type Language struct {
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status     bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted  bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	NativeName string    `gorm:"column:native_name" json:"native_name"`
	Code       string    `gorm:"column:code" json:"code"`
	Locale     string    `gorm:"column:locale" json:"locale"`
	FileID     string    `gorm:"column:file_id" json:"file_id"`

	Flag *File `gorm:"foreignKey:FileID;references:ID" json:"flag"`
}

// TableName Language's table name
func (*Language) TableName() string {
	return TableNameLanguage
}
