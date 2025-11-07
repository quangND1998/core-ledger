package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameOtp = "otp"

// Otp mapped from table <otp>
type Otp struct {
	CreatedAt     time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status        int32      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted     int32      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID            string     `gorm:"column:id;primaryKey" json:"id"`
	Code          string     `gorm:"column:code" json:"code"`
	Type          string     `gorm:"column:type;not null;default:MAIL" json:"type"`
	Email         string     `gorm:"column:email" json:"email"`
	PhoneNumber   string     `gorm:"column:phone_number" json:"phone_number"`
	VerifiedAt    *time.Time `gorm:"column:verified_at" json:"verified_at"`
	ExpiredAt     time.Time  `gorm:"column:expired_at" json:"expired_at"`
	CallingCodeID *string    `gorm:"column:calling_code_id" json:"calling_code_id"`
}

// TableName Otp's table name
func (*Otp) TableName() string {
	return TableNameOtp
}

func (o *Otp) BeforeSave(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return
}

func (o *Otp) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
	return
}
