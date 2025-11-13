package handlers

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// DataProcessHandler xử lý DataProcessJob
type ImportCoaAccountHandler struct {
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
	logger        logger.CustomLogger
	// thêm dependency nếu cần (ví dụ: services, repos)
}

func NewImportCoaAccountHandler(db *gorm.DB, coAccountRepo repo.CoAccountRepo) *ImportCoaAccountHandler {
	return &ImportCoaAccountHandler{
		db:            db,
		coAccountRepo: coAccountRepo,
		logger:        logger.NewSystemLog("ImportCoaAccountHandler"),
	}
}

// NewDataProcessRegistration: provider đăng ký job/handler vào group "queue-registrations"
func NewImportCoaAccountHandlerRegistration(h *ImportCoaAccountHandler) queue.Registration {
	return queue.Registration{
		Type:     "import_coa_account:job",
		Template: &jobs.ImportCoaAccount{},
		Handler:  h,
	}
}

func (h *ImportCoaAccountHandler) Handle(ctx context.Context, j queue.Job) error {
	// kiểu assert về concrete job
	job, ok := j.(*jobs.ImportCoaAccount)
	if !ok {
		return fmt.Errorf("invalid job type, expect *ImportCoaAccount")
	}
	data := job.Data
	h.logger.Info("Data", data.TmpFile)
	if data.TmpFile == "" {
		log.Fatal("File path is empty")
	}
	f, err := excelize.OpenFile(data.TmpFile)

	if err != nil {
		return err
	}
	var accounts []*model.CoaAccount
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		rows, _ := f.GetRows(sheet)
		if len(rows) < 2 {
			continue
		}
		err := h.db.Transaction(func(tx *gorm.DB) error {
			headers := rows[0]
			h.logger.Info("Importing sheet %s with headers %v", sheet, headers)
			// Transaction: nếu lỗi rollback toàn sheet
			for i, row := range rows[1:] {
				rowMap := map[string]string{}
				for j, cell := range row {
					if j < len(headers) {
						rowMap[headers[j]] = strings.TrimSpace(cell)
					}
				}

				name := rowMap["name"]

				code := rowMap["code"]
				Type := rowMap["type"]
				currency := rowMap["currency"]
				parent_code := rowMap["parent_code"]

				var metadataJSON *datatypes.JSON
				if m := rowMap["metadata"]; m != "" {
					var metadataMap map[string]interface{}
					if err := json.Unmarshal([]byte(m), &metadataMap); err != nil {
						h.logger.Warn("Invalid metadata JSON at code %s: %v", code, err)
					} else {
						b, _ := json.Marshal(metadataMap) // convert map -> []byte
						tmp := datatypes.JSON(b)
						metadataJSON = &tmp
					}
				}
				if code == "" {
					h.logger.Warn("Skipping empty code at row %d in sheet %s", i+2, sheet)
					continue
				}

				parentID, _ := h.coAccountRepo.GetParentID(ctx, parent_code)

				// convert string -> *string cho Provider, Network
				var providerPtr, networkPtr *string

				if p := rowMap["provider"]; p != "" {
					provider := p
					providerPtr = &provider
				}

				if n := rowMap["network"]; n != "" {
					network := n
					networkPtr = &network
				}
				account := &model.CoaAccount{
					Code:     code,
					Name:     name,
					Type:     Type,
					Currency: currency,
					ParentID: parentID,
					Provider: providerPtr,
					Network:  networkPtr,
					Metadata: metadataJSON, // TODO: parse JSON nếu cần

				}

				accounts = append(accounts, account)
				// h.logger.Info("Importing co-account %s with data %v", code, accounts)

			}
			h.coAccountRepo.Upsert(accounts, []string{})
			return nil
		})
		if err != nil {
			fmt.Printf("Sheet %s failed: %v\n", sheet, err)
			continue // tiếp tục sheet khác
		}
		fmt.Printf("Sheet %s imported successfully\n", sheet)
	}
	defer os.Remove(data.TmpFile)
	return nil // trả về error để asynq retry nếu cần
}

// Failed: hook được gọi khi job đã hết retry hoặc timeout
func (h *ImportCoaAccountHandler) Failed(ctx context.Context, j queue.Job, err error) {
	// cố gắng assert đúng loại job để log chi tiết
	if job, ok := j.(*jobs.ImportCoaAccount); ok {
		log.Printf("[FAILED] DataProcessJob Type=%s Action=%s Error=%v", job.ProcessType, job.Action, err)
	} else {
		log.Printf("[FAILED] DataProcessJob Error=%v", err)
	}
	// TODO: Có thể ghi log vào DB, tạo transaction_log, hoặc đẩy sang channel cảnh báo...
	_ = h // giữ chỗ nếu sau này cần dùng repo để lưu DB
}
