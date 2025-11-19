package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Entity
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:unique_role_name_guard" json:"name"`
	GuardName string    `gorm:"type:varchar(50);not null;default:'web';uniqueIndex:unique_role_name_guard" json:"guard_name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_has_permissions;joinForeignKey:role_id;joinReferences:permission_id" json:"permissions,omitempty"`
}

func (r *Role) TableName() string {
	return "roles"
}

// GivePermissionTo assigns a permission to role
func (r *Role) GivePermissionTo(db *gorm.DB, permissionName string, guardName string) error {
	if guardName == "" {
		guardName = r.GuardName
	}

	// Find or create permission
	var permission Permission
	err := db.Where("name = ? AND guard_name = ?", permissionName, guardName).FirstOrCreate(&permission, Permission{
		Name:      permissionName,
		GuardName: guardName,
	}).Error
	if err != nil {
		return err
	}

	// Check if already assigned
	var existing RoleHasPermission
	err = db.Where("role_id = ? AND permission_id = ?", r.ID, permission.ID).
		First(&existing).Error
	if err == nil {
		// Already assigned
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Assign permission
	roleHasPermission := RoleHasPermission{
		RoleID:       r.ID,
		PermissionID: permission.ID,
	}
	return db.Create(&roleHasPermission).Error
}

// RevokePermissionTo removes a permission from role
func (r *Role) RevokePermissionTo(db *gorm.DB, permissionName string, guardName string) error {
	if guardName == "" {
		guardName = r.GuardName
	}

	var permission Permission
	err := db.Where("name = ? AND guard_name = ?", permissionName, guardName).First(&permission).Error
	if err != nil {
		return err
	}

	return db.Where("role_id = ? AND permission_id = ?", r.ID, permission.ID).
		Delete(&RoleHasPermission{}).Error
}

// SyncPermissions syncs role permissions (removes all and assigns new ones)
func (r *Role) SyncPermissions(db *gorm.DB, permissionNames []string, guardName string) error {
	if guardName == "" {
		guardName = r.GuardName
	}

	// Remove all existing permissions
	err := db.Where("role_id = ?", r.ID).
		Delete(&RoleHasPermission{}).Error
	if err != nil {
		return err
	}

	// Assign new permissions
	for _, permName := range permissionNames {
		if err := r.GivePermissionTo(db, permName, guardName); err != nil {
			return err
		}
	}

	return nil
}

