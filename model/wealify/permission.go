package model

import (
	"time"
)

const TableNamePermission = "permissions"

// Permission mapped from table <permissions>
type Permission struct {
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status      bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	ParentID    int64     `gorm:"column:parent_id" json:"parent_id"`
	Routes      string    `gorm:"column:routes;not null" json:"routes"`

	Employees []*Employee `gorm:"many2many:employee_permissions;joinForeignKey:PermissionID;joinReferences:EmployeeID" json:"employees"`
}

// TableName Permission's table name
func (*Permission) TableName() string {
	return TableNamePermission
}
