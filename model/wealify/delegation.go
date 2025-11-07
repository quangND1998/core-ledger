package model

import (
	"time"

	"github.com/google/uuid"
)

type DelegationRole string

const (
	DelegationRoleDirector DelegationRole = "DIRECTOR"
	DelegationRoleFinance  DelegationRole = "FINANCE"
	DelegationRoleSeller   DelegationRole = "SELLER"
)

type Delegation struct {
	ID                  uuid.UUID          `gorm:"type:char(36);primaryKey;column:id" json:"id"`
	AccountID           int                `gorm:"column:account_id;index" json:"account_id"`
	DelegationAccountID uuid.UUID          `gorm:"type:char(36);column:delegation_account_id;index" json:"delegation_account_id"`
	DelegationAccount   *DelegationAccount `gorm:"foreignKey:DelegationAccountID;references:ID" json:"delegation_account"`
	StartDate           time.Time          `gorm:"column:start_date;not null" json:"start_date"`
	EndDate             *time.Time         `gorm:"column:end_date" json:"end_date,omitempty"`
	Role                DelegationRole     `gorm:"type:enum('DIRECTOR','FINANCE','SELLER');column:role" json:"role"`
	IsActive            bool               `gorm:"column:is_active;default:true" json:"is_active"`
}

func (Delegation) TableName() string {
	return "delegation"
}
