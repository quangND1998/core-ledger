package handlers

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"fmt"

	"gorm.io/gorm"
)

// LogHandler xử lý LogJob
type LogHandler struct {
	db     *gorm.DB
	logger logger.CustomLogger
}

func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{
		db:     db,
		logger: logger.NewSystemLog("LogHandler"),
	}
}

// NewLogHandlerRegistration: provider đăng ký job/handler vào group "queue-registrations"
func NewLogHandlerRegistration(h *LogHandler) queue.Registration {
	return queue.Registration{
		Type:     "log:create",
		Template: &jobs.LogJob{},
		Handler:  h,
	}
}

func (h *LogHandler) Handle(ctx context.Context, j queue.Job) error {
	// kiểu assert về concrete job
	job, ok := j.(*jobs.LogJob)
	if !ok {
		return fmt.Errorf("invalid job type, expect *LogJob")
	}

	// Tạo log entry
	log := model.Log{
		LoggableID:   job.LoggableID,
		LoggableType: job.LoggableType,
		Action:       job.Action,
		OldValue:     job.OldValue,
		NewValue:     job.NewValue,
		Metadata:     job.Metadata,
		CreatedBy:    job.CreatedBy,
		CreatedAt:    job.CreatedAt,
	}

	// Insert log vào database
	if err := h.db.WithContext(ctx).Table("logs").Create(&log).Error; err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create log for %s:%d: %v", job.LoggableType, job.LoggableID, err))
		return fmt.Errorf("failed to create log: %w", err)
	}

	h.logger.Info(fmt.Sprintf("Successfully created log for %s:%d", job.LoggableType, job.LoggableID))
	return nil
}





