package jobs

import (
	"time"

	"core-ledger/pkg/queue"
)

// DataProcessJob job xử lý dữ liệu theo pattern Laravel
type ImportCoaAccount struct {
	queue.BaseJob
	ProcessType string                 `json:"process_type"`
	Action      string                 `json:"action"`
	Data        DataImportCoaAccount   `json:"data"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

type DataImportCoaAccount struct {
	TmpFile string `json:"tmp_file"`
}

// GetPayload trả về payload của job
func (j *ImportCoaAccount) GetPayload() interface{} {
	return j
}

// GetType trả về loại job
func (j *ImportCoaAccount) GetType() string {
	return "import_coa_account:job"
}

// NewDataProcessJob tạo data process job mới
func NewImportCoaAccount(processType, action string, data DataImportCoaAccount) *ImportCoaAccount {
	return &ImportCoaAccount{
		BaseJob: queue.BaseJob{
			Queue: "default",
			Retry: 1,
		},
		ProcessType: processType,
		Action:      action,
		Data:        data,
	}
}

// SetOptions set processing options
func (j *ImportCoaAccount) SetOptions(options map[string]interface{}) {
	j.Options = options
}

// SetQueue set queue name
func (j *ImportCoaAccount) SetQueue(queue string) {
	j.Queue = queue
}

// SetDelay set delay time
func (j *ImportCoaAccount) SetDelay(delay time.Duration) {
	j.Delay = delay
}

// SetRetry set số lần retry
func (j *ImportCoaAccount) SetRetry(retry int) {
	j.Retry = retry
}
