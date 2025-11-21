package ruleCategory

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"net/http"

	"github.com/gin-gonic/gin"
	// "encoding/json"
)

type RuleCategoryHandler struct {
	logger           logger.CustomLogger
	service          *RuleCateogySerive
	ruleCategoryRepo repo.RuleCategoryRepo
	dispatcher       queue.Dispatcher
}

func NewRuleCategoryHandler(service *RuleCateogySerive, ruleCategoryRepo repo.RuleCategoryRepo, dispatcher queue.Dispatcher) *RuleCategoryHandler {
	return &RuleCategoryHandler{
		logger:           logger.NewSystemLog("CoaAccountHandler"),
		service:          service,
		ruleCategoryRepo: ruleCategoryRepo,
		dispatcher:       dispatcher,
	}
}

func (h *RuleCategoryHandler) List(c *gin.Context) {
	res, err := h.ruleCategoryRepo.List(c)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

// GetCoaAccountRules trả về cấu trúc rules để tạo mã COA account
func (h *RuleCategoryHandler) GetCoaAccountRules(c *gin.Context) {
	res, err := h.service.GetCoaAccountRules(c)
	if err != nil {
		h.logger.Error("Failed to get COA account rules", err)
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}
