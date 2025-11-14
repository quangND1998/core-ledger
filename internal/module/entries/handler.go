package entries

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	// "encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"
)

type EntriesHandler struct {
	logger     logger.CustomLogger
	service    *EntriesService
	entryRepo  repo.EnTriesRepo
	dispatcher queue.Dispatcher
}

func NewEntriesHandler(service *EntriesService, entryRepo repo.EnTriesRepo, dispatcher queue.Dispatcher) *EntriesHandler {
	return &EntriesHandler{
		logger:     logger.NewSystemLog("EntriesHandler"),
		service:    service,
		entryRepo:  entryRepo,
		dispatcher: dispatcher,
	}
}

func (h *EntriesHandler) List(c *gin.Context) {
	// TODO implement me
	q := &dto.ListEntrytFilter{}
	err := c.ShouldBindQuery(&q)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Info("ListEntriesFilter request", q)
	res, err := h.entryRepo.PaginateWithScopes(c, q)
	h.logger.Info("ListEntriesFilter res", res)
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}
