package transactions

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	GetList(c *gin.Context)
}
type transactionHandler struct {
	logger  logger.CustomLogger
	service TransactionService
}

func (h *transactionHandler) GetList(c *gin.Context) {
	ts, err := h.service.listTransaction(c.Request.Context())
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

func NewTransactionHandler(s TransactionService) TransactionHandler {
	return &transactionHandler{
		logger:  logger.NewSystemLog("TransactionHandler"),
		service: s,
	}
}
