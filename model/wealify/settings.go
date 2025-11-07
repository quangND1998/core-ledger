package model

import (
	"core-ledger/model/enum"
	"time"

	"gorm.io/datatypes"
)

const TableNameSetting = "settings"

// Setting mapped from table <settings>
type Setting struct {
	ID        int32           `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Data      datatypes.JSON  `gorm:"column:data;not null" json:"data"`
	CreatedAt time.Time       `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	IsDeleted int32           `gorm:"column:is_deleted;not null" json:"is_deleted"`
	UpdatedBy string          `gorm:"column:updated_by;not null" json:"updated_by"`
	Key       enum.SettingKey `gorm:"column:key" json:"key"`
}

// TableName Setting's table name
func (*Setting) TableName() string {
	return TableNameSetting
}

type AutoChangeApproveTopUp struct {
	Enabled    bool            `json:"enabled"`
	Threshold  uint64          `json:"threshold"`
	PlatformID string          `json:"platform_id"`
	Conditions []RemarkSetting `json:"conditions"`
}

type RemarkSetting struct {
	Type  string `json:"type"` //oneof REGEX,CONTAINS
	Value string `json:"value"`
}
