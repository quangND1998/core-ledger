package coaaccount

import "core-ledger/model/dto"

type ExportRequest struct {
	Select   []CoaAccountExportKey     `json:"select"`
	Query    *dto.ListCoaAccountFilter `json:"query"`
	FileName string                    `json:"file_name"`
}

type ExportResponse struct {
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

type CoaAccountExportKey string

const (
	CoaAccountExportKeyIndex      CoaAccountExportKey = "index"
	CoaAccountExportKeyCode       CoaAccountExportKey = "code"
	CoaAccountExportKeyName       CoaAccountExportKey = "name"
	CoaAccountExportKeyType       CoaAccountExportKey = "type"
	CoaAccountExportKeyParentCode CoaAccountExportKey = "parent_code"
	CoaAccountExportKeyStatus     CoaAccountExportKey = "status"
	CoaAccountExportKeyCurrency   CoaAccountExportKey = "currency"
	CoaAccountExportKeyProvider   CoaAccountExportKey = "provider"

	CoaAccountExportKeyNetwork  CoaAccountExportKey = "network"
	CoaAccountExportKeyTags     CoaAccountExportKey = "tags"
	CoaAccountExportKeyMetadata CoaAccountExportKey = "metadata"

	CoaAccountExportKeyCreatedAt CoaAccountExportKey = "created_at"
	CoaAccountExportKeyUpdatedAt CoaAccountExportKey = "updated_at"
)

func (t CoaAccountExportKey) String() string {
	return string(t)
}

var TransactionExportHeaders = map[CoaAccountExportKey]string{
	"index":       "STT",
	"code":        "code",
	"name":        "name",
	"type":        "type",
	"parent_code": "parent_code",
	"status":      "status",
	"provider":    "provider",
	"currency":    "currency",
	"network":     "network",
	"tags":        "tags",
	"metadata":    "metadata",
	"created_at":  "Ngày tạo",
	"updated_at":  "Ngày cập nhật",
}
