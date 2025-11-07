package model

import (
	"core-ledger/model/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const TableNameNotification = "notifications"

// Notification mapped from table <notifications>
type Notification struct {
	CreatedAt           time.Time               `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt           time.Time               `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status              bool                    `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted           bool                    `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID                  string                  `gorm:"column:id;primaryKey" json:"id"`
	TargetID            *string                 `gorm:"column:target_id" json:"target_id"`
	Title               string                  `gorm:"column:title" json:"title"`
	Subject             string                  `gorm:"column:subject" json:"subject"`
	Description         string                  `gorm:"column:description" json:"description"`
	Link                string                  `gorm:"column:link" json:"link"`
	OldNotificationType string                  `gorm:"column:old_notification_type" json:"old_notification_type"`
	NotificationGroup   enum.NotificationGroup  `gorm:"column:notification_group;not null" json:"notification_group"`
	NotificationStatus  enum.NotificationStatus `gorm:"column:notification_status;not null" json:"notification_status"`
	CustomerID          *int64                  `gorm:"column:customer_id" json:"customer_id"`
	EmployeeID          *int64                  `gorm:"column:employee_id" json:"employee_id"`
	FileID              *string                 `gorm:"column:file_id" json:"file_id"`
	Type                enum.NotificationType   `gorm:"column:type;not null" json:"type"`
	Data                *datatypes.JSON         `gorm:"column:data" json:"data"`
}

// TableName Notification's table name
func (*Notification) TableName() string {
	return TableNameNotification
}

func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return
}

func (n *Notification) BeforeSave(tx *gorm.DB) (err error) {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return
}
