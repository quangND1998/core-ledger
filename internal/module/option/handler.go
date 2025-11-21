package option

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"

	// "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OptionHandler struct {
	db         *gorm.DB
	logger     logger.CustomLogger
	service    *OptionsSerive
	dispatcher queue.Dispatcher
}

// NewOptionHandler - DEPRECATED: Handler này đã không còn sử dụng model cũ
// Có thể xóa hoàn toàn nếu không cần thiết
func NewOptionHandler(db *gorm.DB, service *OptionsSerive, dispatcher queue.Dispatcher) *OptionHandler {
	return &OptionHandler{
		db:         db,
		logger:     logger.NewSystemLog("OptionHandler"),
		service:    service,
		dispatcher: dispatcher,
	}
}

// func (h *OptionHandler) GetRuleTypes(c *gin.Context) {
// 	data, err := h.service.GetRuleTypes(c.Request.Context())
// 	if err != nil {
// 		h.logger.Error("failed to get transaction list: %v", err)
// 		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	ginhp.RespondOK(c, data)
// }

// func (h *OptionHandler) GetRuleGroups(c *gin.Context) {

// 	typeId, err := utils.ParseIntIdParam(c.Param("typeId"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Invalid SystemPaymentID"})
// 		return
// 	}
// 	data, err := h.service.GetRuleGroups(c.Request.Context(), typeId)
// 	if err != nil {
// 		h.logger.Error("failed to get transaction list: %v", err)
// 		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	ginhp.RespondOK(c, data)
// }
// func (h *OptionHandler) GetRuleOptionSteps(c *gin.Context) {
// 	groupId := c.Param("groupId")

// 	var steps []model.AccountRuleOptionStep
// 	h.db.Where("option_id = ?", groupId).
// 		Order("step_order ASC").
// 		Find(&steps)

// 	var resp []dto.StepResponse

// 	for _, step := range steps {
// 		item := dto.StepResponse{
// 			StepID:    step.ID,
// 			StepOrder: step.StepOrder,
// 			Label:     step.InputLabel,
// 		}

// 		// CASE: dropdown
// 		if step.CategoryID != nil {
// 			var category model.RuleCategory
// 			h.db.First(&category, *step.CategoryID)

// 			var values []model.RuleValue
// 			h.db.Where("category_id = ? AND is_delete = false", *step.CategoryID).
// 				Order("sort_order ASC").
// 				Find(&values)

// 			item.Type = "dropdown"
// 			item.CategoryCode = category.Code
// 			item.Values = values
// 		}

// 		// CASE: input
// 		if step.CategoryID == nil && *step.InputCode != "" {
// 			item.Type = "input"
// 			item.InputType = step.InputType
// 		}

// 		resp = append(resp, item)
// 	}

// 	c.JSON(200, resp)
// }

// GetRuleOptionTree - DEPRECATED: Sử dụng cấu trúc mới coa_account_rule_types
// func (h *OptionHandler) GetRuleOptionTree(c *gin.Context) {
// 	var options []model.AccountRuleOption
// 	h.db.Order("parent_option_id ASC").Find(&options)
// 	...
// }

// GetFullRuleData - DEPRECATED: Sử dụng endpoint /rule-category/coa-rules thay thế
// func (h *OptionHandler) GetFullRuleData(c *gin.Context) {
// 	...
// }

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
