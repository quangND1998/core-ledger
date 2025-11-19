package ruleValue

import (
	"core-ledger/internal/module/validate"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "encoding/json"
)

type RuleValueHandler struct {
	db     *gorm.DB
	logger logger.CustomLogger
	service          *RuleCateogySerive
	ruleValueRepo    repo.RuleValueRepo
	ruleCategoryRepo repo.RuleCategoryRepo
	dispatcher       queue.Dispatcher
}

func NewRuleValueHandler(db *gorm.DB, service *RuleCateogySerive, ruleValueRepo repo.RuleValueRepo, ruleCategoryRepo repo.RuleCategoryRepo, dispatcher queue.Dispatcher) *RuleValueHandler {
	return &RuleValueHandler{
		db:               db,
		logger:           logger.NewSystemLog("CoaAccountHandler"),
		service:          service,
		ruleValueRepo:    ruleValueRepo,
		ruleCategoryRepo: ruleCategoryRepo,
		dispatcher:       dispatcher,
	}
}

func (h *RuleValueHandler) List(c *gin.Context) {
	var req dto.FilterRuleValueRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Error while binding query", err)
		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid query", out)

		return
	}

	res, err := h.ruleValueRepo.List(c, &req)
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *RuleValueHandler) Create(c *gin.Context) {
	var req dto.SaveRuleValueRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)

		return
	}

	if len(req.Data) == 0 {
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", map[string]string{"Data": "cannot be empty"})
		return
	}
	if errs := h.ValidateUniqueRuleCode(req.Data, h.db); errs != nil {
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Duplicate values", errs)
		return
	}
	_, err := h.ruleCategoryRepo.GetByID(c, int64(req.CategoryID))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Category Rule not found")
		return
	}

	var toUpsert []*model.RuleValue

	for _, v := range req.Data {
		// nếu is_delete = false → upsert
		toUpsert = append(toUpsert, &model.RuleValue{
			ID:         v.ID,
			CategoryID: req.CategoryID,
			Value:      v.Value,
			Name:       v.Name,
			IsDelete:   *v.IsDelete,
		})

	}

	// Upsert các bản ghi còn lại
	if len(toUpsert) > 0 {
		if err := h.ruleValueRepo.Upsert(toUpsert, []string{}); err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to upsert records")
			return
		}
	}

	ginhp.RespondOK(c, "Update rule value successfully")
}

func (h *RuleValueHandler) ValidateUniqueRuleCode(data []*dto.RuleValueRequest, db *gorm.DB) map[string]string {
	errors := map[string]string{}
	for i, item := range data {
		if item.IsDelete != nil && *item.IsDelete {
			// Bỏ qua những item đánh dấu xóa
			continue
		}

		var count int64
		query := db.Model(&model.RuleValue{}).Where("value = ? AND is_delete = false", item.Value)
		if item.ID != 0 {
			query = query.Where("id != ?", item.ID)
		}
		if err := query.Count(&count).Error; err != nil {
			errors[fmt.Sprintf("data.%d.value", i)] = "Database error"
			continue
		}
		if count > 0 {
			errors[fmt.Sprintf("data.%d.value", i)] = "Value already exists"
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
