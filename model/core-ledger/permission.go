package model

import (
	"time"
)

type Permission struct {
	Entity
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:unique_permission_name_guard" json:"name"`
	GuardName string    `gorm:"type:varchar(50);not null;default:'web';uniqueIndex:unique_permission_name_guard" json:"guard_name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	Roles []Role `gorm:"many2many:role_has_permissions;joinForeignKey:permission_id;joinReferences:role_id" json:"roles,omitempty"`
}

func (p *Permission) TableName() string {
	return "permissions"
}

