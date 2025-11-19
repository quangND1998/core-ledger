package model

type RoleHasPermission struct {
	PermissionID uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	RoleID       uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`

	// Relations
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
	Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (r *RoleHasPermission) TableName() string {
	return "role_has_permissions"
}

