package model

type ModelHasRole struct {
	RoleID   uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`
	ModelID  uint64 `gorm:"primaryKey;column:model_id" json:"model_id"`
	ModelType string `gorm:"primaryKey;type:varchar(255);column:model_type" json:"model_type"`

	// Relations
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (m *ModelHasRole) TableName() string {
	return "model_has_roles"
}

