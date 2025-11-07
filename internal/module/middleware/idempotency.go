package middleware

import (
	"bytes"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logging"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdempotencyGetter interface {
	Get(*gin.Context) (string, error)
}

// IdempotencyMiddleware creates a new Gin middleware for handling idempotency.
func IdempotencyMiddleware(db *gorm.DB, getIdempotencyFn func(c *gin.Context) (string, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logging.From(c).Named("IdempotencyMiddleware")
		key, err := getIdempotencyFn(c)
		if key == "" {
			logger.Errorw("empty idempotency key", "err", err)
			// Or just let it pass if you want idempotency to be optional
			c.Next()
			return
		}
		logger = logger.With("idempotency_key", key)

		// Read and hash the request body. We need to do this carefully
		// so we can "put it back" for the actual handler.
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body
		requestHash := hashBytes(bodyBytes)

		// Get the user ID from the context (set by a previous auth middleware)
		employee := ginhp.GetEmployeeReq(c)

		employeeID := employee.ID

		var record model.IdempotencyRecord

		// Use a transaction to handle locking and prevent race conditions
		err = db.Transaction(func(tx *gorm.DB) error {
			// Lock the row to prevent another request with the same key from processing
			err := tx.Set("gorm:query_option", "FOR UPDATE").Where("idempotency_key = ?", key).First(&record).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				// Real database error
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// --- CASE: KEY IS NEW ---
				newRecord := model.IdempotencyRecord{
					IdempotencyKey: key,
					EmployeeID:     employeeID,
					RequestHash:    requestHash,
					Status:         model.IdempotencyInProgress,
					ExpiresAt:      time.Now().Add(5 * time.Minute),
				}
				if err := tx.Create(&newRecord).Error; err != nil {
					return err
				}
				// Mark context to indicate we should save the response later
				c.Set("is_idempotent_new", true)
				return nil // Commit the "IN_PROGRESS" state
			}

			// --- CASE: KEY EXISTS ---
			if record.RequestHash != requestHash {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Idempotency key is being reused with a different request body."})
				return errors.New("idempotency conflict") // Rollback
			}

			if record.Status == "IN_PROGRESS" {
				ginhp.RespondError(c, http.StatusConflict, "Request đang được xử lý")
				c.Abort()
				return errors.New("idempotency conflict") // Rollback
			}

			if record.Status == "COMPLETED" || record.Status == "FAILED" {
				// Return the saved response and stop processing
				c.Data(record.ResponseCode, "application/json", record.ResponseBody)
				c.Abort()
				return errors.New("idempotency completed") // Rollback
			}
			return nil
		})

		if err != nil {
			// The transaction errored or was intentionally aborted (e.g., key found),
			// the response has already been sent by c.AbortWith...
			return
		}

		// ---- Let the actual handler run ----
		// We need a custom response writer to capture the body
		responseWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next() // <<<<<<< HANDLER EXECUTES HERE >>>>>>>

		// ---- After the handler ----
		if isNew, exists := c.Get("is_idempotent_new"); exists && isNew.(bool) {
			// Update the record with the final result
			finalStatus := model.IdempotencyCompleted
			if c.Writer.Status() >= 400 {
				finalStatus = model.IdempotencyFailed
			}
			db.Model((*model.IdempotencyRecord)(nil)).Where("idempotency_key = ?", key).Updates(model.IdempotencyRecord{
				Status:       finalStatus,
				ResponseCode: c.Writer.Status(),
				ResponseBody: responseWriter.body.Bytes(),
			})
		}
	}
}

func hashBytes(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

// Helper to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
