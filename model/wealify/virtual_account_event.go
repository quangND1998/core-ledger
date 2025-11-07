package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameVirtualAccountEvent = "virtual_account_events"

type VirtualAccountEvent struct {
	ID                  string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt           time.Time `gorm:"created_at" json:"created_at,omitempty"`
	EventType           string    `gorm:"type:varchar(20);default:'CREATE'" json:"event_type,omitempty"`
	OldStatus           *VAStatus `gorm:"type:varchar(20)" json:"old_status,omitempty"`
	NewStatus           *VAStatus `gorm:"type:varchar(20)" json:"new_status,omitempty"`
	VirtualAccountID    int64     `gorm:"column:virtual_account_id" json:"virtual_account_id,omitempty"`
	CreatedByCustomerID *int64    `gorm:"column:created_by_customer_id" json:"created_by_customer_id,omitempty"`
	CreatedByEmployeeID *int64    `gorm:"column:created_by_employee_id" json:"created_by_employee_id,omitempty"`
	CreatedByType       string    `gorm:"column:created_by_type" json:"created_by_type,omitempty"`

	VirtualAccount    *VirtualAccount `gorm:"foreignkey:VirtualAccountID;references:ID" json:"virtual_account,omitempty"`
	CreatedByCustomer *Customer       `gorm:"foreignkey:CreatedByCustomerID;references:ID" json:"created_by_customer,omitempty"`
	CreatedByEmployee *Employee       `gorm:"foreignkey:CreatedByEmployeeID;references:ID" json:"created_by_employee,omitempty"`
}

// TableName VirtualAccount's table name
func (*VirtualAccountEvent) TableName() string {
	return TableNameVirtualAccountEvent
}

func (v *VirtualAccountEvent) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return nil
}

func (v *VirtualAccountEvent) BeforeSave(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return nil
}
