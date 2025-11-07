package model

import (
	"time"
)

const TableNameFile = "files"

// File mapped from table <files>
type File struct {
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status     bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted  bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID         string    `gorm:"column:id;primaryKey" json:"id"`
	Driver     string    `gorm:"column:driver;not null;default:AWS" json:"driver"`
	Name       string    `gorm:"column:name" json:"name"`
	Token      string    `gorm:"column:token" json:"token"`
	CustomerID int64     `gorm:"column:customer_id" json:"customer_id"`
	EmployeeID int64     `gorm:"column:employee_id" json:"employee_id"`
}

// TableName File's table name
func (*File) TableName() string {
	return TableNameFile
}
