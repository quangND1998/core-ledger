package model

import (
	"time"

	"github.com/google/uuid"
)

type DelegationAccountStatus string

const (
	DelegationAccountStatusInvited  DelegationAccountStatus = "INVITED"
	DelegationAccountStatusActive   DelegationAccountStatus = "ACTIVE"
	DelegationAccountStatusInactive DelegationAccountStatus = "INACTIVE"
	// thêm các trạng thái khác nếu có
)

type DelegationAccount struct {
	ID        uuid.UUID               `gorm:"type:char(36);primaryKey;column:id" json:"id"`
	Username  string                  `gorm:"type:varchar(100);uniqueIndex;column:username" json:"username"`
	Email     string                  `gorm:"type:varchar(255);uniqueIndex;column:email" json:"email"`
	Password  string                  `gorm:"type:varchar(255);column:password" json:"password"`
	Status    DelegationAccountStatus `gorm:"type:enum('INVITED','ACTIVE','INACTIVE');default:INVITED;column:status" json:"status"`
	CreatedAt time.Time               `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time               `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (DelegationAccount) TableName() string {
	return "delegation_account"
}
