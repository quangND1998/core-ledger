package option

import (
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OptionHandler struct {
	db         *gorm.DB
	logger     logger.CustomLogger
	service    *OptionsSerive
	dispatcher queue.Dispatcher
}

func NewOptionHandler(db *gorm.DB, service *OptionsSerive, dispatcher queue.Dispatcher) *OptionHandler {
	return &OptionHandler{
		db:         db,
		logger:     logger.NewSystemLog("CoaAccountHandler"),
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

func (h *OptionHandler) GetRuleOptionTree(c *gin.Context) {
	var options []model.AccountRuleOption
	h.db.Order("parent_option_id ASC").Find(&options)

	// Index by ID (uint64)
	m := make(map[uint64]*dto.AccountRuleOptionTree)

	for _, o := range options {
		m[o.ID] = &dto.AccountRuleOptionTree{
			ID:       o.ID,
			Code:     o.Code,
			Name:     o.Name,
			LayerID:  o.LayerID,
			ParentID: o.ParentOptionID,
			Children: []*dto.AccountRuleOptionTree{},
		}
	}

	var roots []*dto.AccountRuleOptionTree

	for _, node := range m {
		if node.ParentID == nil {
			roots = append(roots, node)
		} else {
			parent := m[*node.ParentID] // now uint64
			parent.Children = append(parent.Children, node)
		}
	}

	c.JSON(200, roots)
}

func (h *OptionHandler) GetFullRuleData(c *gin.Context) {
	var allOptions []model.AccountRuleOption
	var allSteps []model.AccountRuleOptionStep
	var allCategories []model.RuleCategory
	var allValues []model.RuleValue

	// Load toàn bộ data
	h.db.Order("sort_order ASC").Find(&allOptions)
	h.db.Order("step_order ASC").Find(&allSteps)
	h.db.Find(&allCategories)
	h.db.Where("is_delete = false").Order("sort_order ASC").Find(&allValues)

	// --------- Map index ---------
	categoryMap := make(map[uint64]model.RuleCategory)
	for _, ctg := range allCategories {
		categoryMap[uint64(ctg.ID)] = ctg
	}

	valueMap := make(map[uint64][]dto.RuleValueResp)
	for _, v := range allValues {
		valueMap[uint64(v.CategoryID)] = append(valueMap[uint64(v.CategoryID)], dto.RuleValueResp{
			ID:    uint64(v.ID),
			Name:  v.Name,
			Value: v.Value,
		})
	}

	// --------- Group option theo parent ---------
	optionMap := make(map[uint64][]model.AccountRuleOption) // parent id -> children

	for _, opt := range allOptions {
		parentID := uint64(0)
		if opt.ParentOptionID != nil {
			parentID = *opt.ParentOptionID
		}
		optionMap[parentID] = append(optionMap[parentID], opt)
	}

	// --------- Build response ---------
	var result []dto.RuleTypeResp

	// TYPE = parent_option_id IS NULL (tức parentID = 0)
	for _, t := range optionMap[0] {
		typeResp := dto.RuleTypeResp{
			ID:   t.ID,
			Code: t.Code,
			Name: t.Name,
		}

		// GROUP của TYPE này
		for _, g := range optionMap[t.ID] {
			groupResp := dto.RuleGroupResp{
				ID:        g.ID,
				Code:      g.Code,
				Name:      g.Name,
				InputType: g.InputType,
			}

			// STEP của group
			var steps []dto.RuleStepResp
			for _, s := range allSteps {
				if s.OptionID == g.ID {
					stepResp := dto.RuleStepResp{
						StepID:    s.ID,
						StepOrder: s.StepOrder,
						Label:     valueOrEmpty(s.InputLabel),
						InputCode: s.InputCode,
						InputType: s.InputType,
					}

					// dropdown
					if s.CategoryID != nil {
						cat := categoryMap[*s.CategoryID]
						stepResp.Type = s.InputType
						stepResp.CategoryCode = cat.Code
						stepResp.Values = valueMap[uint64(cat.ID)]
					}

					// input
					if s.CategoryID == nil && s.InputCode != nil {
						stepResp.Type = s.InputType
					}

					steps = append(steps, stepResp)
				}
			}

			groupResp.Steps = steps
			typeResp.Groups = append(typeResp.Groups, groupResp)
		}

		result = append(result, typeResp)
	}

	c.JSON(200, result)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
