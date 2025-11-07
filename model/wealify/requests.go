package model

import (
	"time"
)

const TableNameRequest = "requests"

// Request mapped from table <requests>
type Request struct {
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status      bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	TargetID    string    `gorm:"column:target_id;comment:SystemPaymentID of record, employee want to customer provider more info." json:"target_id"`
	Title       string    `gorm:"column:title;not null" json:"title"`
	Subject     string    `gorm:"column:subject;not null" json:"subject"`
	Description string    `gorm:"column:description;not null" json:"description"`
	RequestType string    `gorm:"column:request_type;not null" json:"request_type"`
	EmployeeID  int32     `gorm:"column:employee_id" json:"employee_id"`
	CustomerID  int32     `gorm:"column:customer_id" json:"customer_id"`
}

// TableName Request's table name
func (*Request) TableName() string {
	return TableNameRequest
}
