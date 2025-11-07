package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameSession = "sessions"

// Session mapped from table <sessions>
type Session struct {
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status         int64     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted      int64     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID             string    `gorm:"column:id;primaryKey" json:"id"`
	Secret         string    `gorm:"column:secret" json:"secret"`
	Token          string    `gorm:"column:token" json:"token"`
	Type           string    `gorm:"column:type;not null" json:"type"`
	ExpiredAt      time.Time `gorm:"column:expired_at" json:"expired_at"`
	CustomerID     *int64    `gorm:"column:customer_id" json:"customer_id"`
	EmployeeID     *int64    `gorm:"column:employee_id" json:"employee_id"`
	IP             string    `gorm:"column:ip" json:"ip"`
	BrowserID      string    `gorm:"column:browser_id" json:"browser_id"`
	DeviceFamily   string    `gorm:"column:device_family" json:"device_family"`
	DeviceVersion  string    `gorm:"column:device_version" json:"device_version"`
	BrowserFamily  string    `gorm:"column:browser_family" json:"browser_family"`
	BrowserVersion string    `gorm:"column:browser_version" json:"browser_version"`
	OsFamily       string    `gorm:"column:os_family" json:"os_family"`
	OsMajor        string    `gorm:"column:os_major" json:"os_major"`
	OsMinor        string    `gorm:"column:os_minor" json:"os_minor"`
	DeviceID       string    `gorm:"column:device_id" json:"device_id"`
}

// TableName Session's table name
func (*Session) TableName() string {
	return TableNameSession
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.NewString()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Session) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

type ActiveSession struct {
	ID        string
	CreatedAt time.Time
}
