package model

import (
	"time"
)

const TableNameToken = "tokens"

// Token mapped from table <tokens>
type Token struct {
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status     bool       `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted  int32      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID         string     `gorm:"column:id;primaryKey" json:"id"`
	Token      string     `gorm:"column:token;not null" json:"token"`
	VerifiedAt *time.Time `gorm:"column:verified_at" json:"verified_at"`
	ExpiredAt  time.Time  `gorm:"column:expired_at" json:"expired_at"`
	CustomerID *int64     `gorm:"column:customer_id" json:"customer_id"`
	EmployeeID *int64     `gorm:"column:employee_id" json:"employee_id"`
}

// TableName Token's table name
func (*Token) TableName() string {
	return TableNameToken
}
