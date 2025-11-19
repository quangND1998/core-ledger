package model

type ModelHasPermission struct {
	PermissionID uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"`
	ModelID      uint64 `gorm:"primaryKey;column:model_id" json:"model_id"`
	ModelType    string `gorm:"primaryKey;type:varchar(255);column:model_type" json:"model_type"`

	// Relations
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

func (m *ModelHasPermission) TableName() string {
	return "model_has_permissions"
}

