package coaaccount

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CoaAccountHandler struct {
	logger        logger.CustomLogger
	service       *CoaAccountService
	coAccountRepo repo.CoAccountRepo
	dispatcher    queue.Dispatcher
}

func NewCoaAccountHandler(service *CoaAccountService, coAccountRepo repo.CoAccountRepo, dispatcher queue.Dispatcher) *CoaAccountHandler {
	return &CoaAccountHandler{
		logger:        logger.NewSystemLog("CoaAccountHandler"),
		service:       service,
		coAccountRepo: coAccountRepo,
		dispatcher:    dispatcher,
	}
}

func (h *CoaAccountHandler) List(c *gin.Context) {
	// TODO implement me
	var q dto.ListCoaAccountFilter
	err := c.ShouldBindQuery(&q)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("ListCoaAccountFilter request")
	res, err := h.coAccountRepo.Paginate(&q)
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}
