package dto

import (
	model "core-ledger/model/wealify"
	"time"

	"gorm.io/datatypes"
)

type QuerySettingDto struct {
	Key string `json:"key" form:"key" binding:"required" example:"va_limit"`
}

type UpdateSettingDto struct {
	Data datatypes.JSON `json:"data" binding:"required"`
}

type UpdateVirtualAccountLimitDto struct {
	Standard *int64 `json:"standard,omitempty" example:"1000000"`
	Silver   *int64 `json:"silver,omitempty" example:"5000000"`
	Gold     *int64 `json:"gold,omitempty" example:"10000000"`
	Diamond  *int64 `json:"diamond,omitempty" example:"50000000"`
}

type EditCfgAutoProcessToApproveRequest struct {
	Configs []model.AutoChangeApproveTopUp `json:"configs"`
}

type RemarkSetting struct {
	Type  string `json:"type" binding:"oneof=REGEX CONTAINS" example:"CONTAINS"`
	Value string `json:"value" binding:"required" example:"topup"`
}

type SettingResponse struct {
	ID        int64          `json:"id"`
	Key       string         `json:"key"`
	Data      datatypes.JSON `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type SettingListResponse struct {
	Items  []SettingResponse `json:"items"`
	Total  int64             `json:"total"`
	Page   int               `json:"page"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

type AutoProcessConfigResponse struct {
	Enabled    bool            `json:"enabled"`
	Threshold  uint64          `json:"threshold"`
	Conditions []RemarkSetting `json:"conditions"`
}
type UpdateAutoApproveTopUpCustomerRequest struct {
	NewIDs    []int64 `json:"new_ids,omitempty"`
	DeleteIDs []int64 `json:"delete_ids,omitempty"`
}
