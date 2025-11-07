package model

import (
	"time"

	"github.com/google/uuid"
)

type DelegationScope struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	VirtualAccountID    int64  `gorm:"column:virtual_account_id;not null" json:"virtual_account_id"`
	DelegationAccountID string `gorm:"type:varchar(255);column:delegation_account_id;not null" json:"delegation_account_id"`

	VirtualAccount    *VirtualAccount    `gorm:"foreignKey:VirtualAccountID;references:ID" json:"virtual_account"`
	DelegationAccount *DelegationAccount `gorm:"foreignKey:DelegationAccountID;references:ID;constraint:onDelete:CASCADE" json:"delegation_account"`
}

func (DelegationScope) TableName() string {
	return "virtual_account_delegation"
}
