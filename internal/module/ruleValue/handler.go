package ruleValue

import (
	"core-ledger/internal/module/validate"
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	// "encoding/json"
)

type RuleValueHandler struct {
	logger        logger.CustomLogger
	service       *RuleCateogySerive
	ruleValueRepo repo.RuleValueRepo
	dispatcher    queue.Dispatcher
}

func NewRuleValueHandler(service *RuleCateogySerive, ruleValueRepo repo.RuleValueRepo, dispatcher queue.Dispatcher) *RuleValueHandler {
	return &RuleValueHandler{
		logger:        logger.NewSystemLog("CoaAccountHandler"),
		service:       service,
		ruleValueRepo: ruleValueRepo,
		dispatcher:    dispatcher,
	}
}

func (h *RuleValueHandler) List(c *gin.Context) {
	res, err := h.ruleValueRepo.List(c)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *RuleValueHandler) Create(c *gin.Context) {
	var req SaveRuleValueRequest // NOT a pointer

	// Bind JSON vào struct

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if ok := errors.As(err, &ve); ok {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = validate.ValidationErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}

		// Other errors
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}
