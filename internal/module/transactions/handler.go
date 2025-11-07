package transactions

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	logger  logger.CustomLogger
	service *transactionService
}

func (h *TransactionHandler) GetList(c *gin.Context) {
	ts, err := h.service.ListTransaction(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get transaction list: %v", err)
		c.JSON(500, gin.H{
			"error": "internal server error ",
			"err":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: ts,
	})
}

func NewTransactionHandler(service *transactionService) *TransactionHandler {
	return &TransactionHandler{
		logger:  logger.NewSystemLog("TransactionHandler"),
		service: service,
	}
}
