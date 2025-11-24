package coaaccount

import (
	"bytes"
	"core-ledger/internal/module/middleware"
	"core-ledger/internal/module/validate"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"
	"core-ledger/pkg/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestCoaAccountHandler struct {
	logger            logger.CustomLogger
	service           *RequestCoaAccountService
	coAccountRepo     repo.CoAccountRepo
	reqCoaAccountRepo repo.ReqCoaAccountRepo
	db                *gorm.DB
}

func NewRequestCoaAccountHandler(service *RequestCoaAccountService, db *gorm.DB, coAccountRepo repo.CoAccountRepo, reqCoaAccountRepo repo.ReqCoaAccountRepo) *RequestCoaAccountHandler {
	return &RequestCoaAccountHandler{
		logger:            logger.NewSystemLog("RequestCoaAccountHandler"),
		service:           service,
		coAccountRepo:     coAccountRepo,
		reqCoaAccountRepo: reqCoaAccountRepo,
		db:                db,
	}
}

// GetList godoc
// @Summary Get list of COA account requests
// @Description Get paginated list of COA account requests with filters
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param request_type query string false "Request type (CREATE, EDIT)"
// @Param request_status query string false "Request status (PENDING, APPROVED, REJECTED)"
// @Param maker_id query int false "Maker ID"
// @Param checker_id query int false "Checker ID"
// @Param coa_account_id query int false "COA Account ID"
// @Param search query string false "Search by account_no, name, code"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.PreResponse
// @Failure 400 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts [get]
// GetList handles GET /request-coa-accounts
func (h *RequestCoaAccountHandler) GetList(c *gin.Context) {
	var filter dto.ListRequestCoaAccountFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetList(c, &filter)
	total_pending_request, err := h.reqCoaAccountRepo.CountByFields(c, map[string]interface{}{
		"request_status": model.RequestStatusPending,
	})
	if err != nil {
		h.logger.Error("Failed to get request list", err)
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	res.TotalPendingRequest = total_pending_request
	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

// GetDetail godoc
// @Summary Get COA account request detail
// @Description Get detailed information of a specific COA account request
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} dto.PreResponse{data=dto.RequestCoaAccountResponse}
// @Failure 400 {object} dto.PreResponse
// @Failure 404 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts/{id} [get]
// GetDetail handles GET /request-coa-accounts/:id
func (h *RequestCoaAccountHandler) GetDetail(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	res, err := h.service.GetByID(c, uint64(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusNotFound, "Request not found")
			return
		}
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

// Create godoc
// @Summary Create COA account request
// @Description Create a new COA account request (CREATE or EDIT type). Request will be created with PENDING status and requires approval.
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param request body dto.RequestCoaAccountCreateRequestWithValidation true "Create request (request_type=CREATE)" example({"request_type":"CREATE","account_data":{"code":"ASSET","account_no":"ASSET:VND.Cobo.Details","name":"Asset Account","type":"ASSET","currency":"VND","status":"ACTIVE"}})
// @Param request body dto.RequestCoaAccountEditRequestWithValidation true "Edit request (request_type=EDIT)" example({"request_type":"EDIT","account_data":{"account_id":1,"account_no":"ASSET:VND.Cobo.NewDetails","name":"Updated Account Name","status":"ACTIVE"}})
// @Success 200 {object} ginhp.Response
// @Failure 400 {object} dto.PreResponse
// @Failure 401 {object} dto.PreResponse
// @Failure 422 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts [post]
// Create handles POST /request-coa-accounts
// Theo luồng: Check duplicate → Tạo request với status PENDING
func (h *RequestCoaAccountHandler) Create(c *gin.Context) {
	// Lấy userID từ context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Read request body to determine request_type first
	bodyBytes, err := c.GetRawData()
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Parse request_type from JSON
	var temp struct {
		RequestType model.RequestType `json:"request_type"`
	}
	if err := json.Unmarshal(bodyBytes, &temp); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Restore body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var request *model.RequestCoaAccount
	var accountNo string

	// Bind to appropriate struct based on request_type
	if temp.RequestType == model.RequestTypeCreate {
		var req dto.RequestCoaAccountCreateRequestWithValidation
		if err := c.ShouldBindJSON(&req); err != nil {
			out := validate.FormatErrorMessage(req, err)
			ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
			return
		}

		accountNo = req.AccountData.AccountNo

		// Check duplicate
		duplicateInfo, err := h.service.CheckDuplicate(c, accountNo, model.RequestTypeCreate)
		if err != nil {
			h.logger.Error("Failed to check duplicate", err)
			ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi kiểm tra trùng lặp")
			return
		}

		if duplicateInfo.IsDuplicate {
			c.JSON(http.StatusBadRequest, dto.PreResponse{
				Data: map[string]interface{}{
					"is_duplicate": true,
					"message":      duplicateInfo.Message,
					"duplicate_in": duplicateInfo.DuplicateIn,
				},
			})
			return
		}
		//check duplicate code
		exists, err := h.coAccountRepo.Exists(c, map[string]any{
			"code":   req.AccountData.Code,
			"status": "ACTIVE",
		})
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra")
			return
		}
		if exists {
			ginhp.RespondError(c, http.StatusBadRequest, "Mã tài khoản đã tồn tại")
			return
		}
		// Convert to model
		request, err = req.ToModel(uint64(userID))
		if err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
	} else if temp.RequestType == model.RequestTypeEdit {
		var req dto.RequestCoaAccountEditRequestWithValidation
		if err := c.ShouldBindJSON(&req); err != nil {
			out := validate.FormatErrorMessage(req, err)
			ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
			return
		}

		accountNo = req.AccountData.AccountNo

		// Kiểm tra account có tồn tại không và lấy account hiện tại
		existingAccount, err := h.coAccountRepo.GetOneByFields(c, map[string]interface{}{
			"id": req.AccountData.AccountId,
		})
		if err != nil {
			ginhp.ReturnBadRequestError(c, err)
			return
		}

		// Chỉ check duplicate nếu account_no thay đổi (khác với account_no hiện tại)
		if existingAccount.AccountNo != accountNo {
			// Check duplicate (loại trừ account hiện tại đang được edit)

			duplicateInfo, err := h.service.CheckDuplicate(c, accountNo, model.RequestTypeEdit)
			if err != nil {
				h.logger.Error("Failed to check duplicate", err)
				ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi kiểm tra trùng lặp")
				return
			}

			if duplicateInfo.IsDuplicate {
				c.JSON(http.StatusBadRequest, dto.PreResponse{
					Data: map[string]interface{}{
						"is_duplicate": true,
						"message":      duplicateInfo.Message,
						"duplicate_in": duplicateInfo.DuplicateIn,
					},
				})
				return
			}
		}
		// Nếu account_no không thay đổi (giữ nguyên), không cần check duplicate

		// Convert to model
		request, err = req.ToModel(uint64(userID))
		if err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid request_type. Must be CREATE or EDIT")
		return
	}

	if err := h.service.Create(c, request); err != nil {
		h.logger.Error("Failed to create request", err)
		ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi tạo request")
		return
	}

	ginhp.RespondOK(c, "Tạo request thành công. Đang chờ phê duyệt.")
}

// Update godoc
// @Summary Update COA account request
// @Description Update a COA account request. Only allowed for PENDING requests with type CREATE or REJECTED requests.
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body dto.RequestCoaAccountCreateRequestWithValidation true "Update CREATE request (for PENDING requests)"
// @Param request body dto.RequestCoaAccountUpdateRequestWithValidation true "Update REJECTED request (only account_no, status, description)"
// @Success 200 {object} ginhp.Response
// @Failure 400 {object} dto.PreResponse
// @Failure 404 {object} dto.PreResponse
// @Failure 422 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts/{id} [put]
// Update handles PUT /request-coa-accounts/:id
// Cho phép update request có status PENDING với request_type = CREATE
// Hoặc request có status REJECTED
func (h *RequestCoaAccountHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Lấy request hiện tại
	request, err := h.service.GetRequestByID(c, uint64(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusNotFound, "Request not found")
			return
		}
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Chỉ cho phép update request PENDING với type CREATE hoặc request REJECTED
	if request.RequestStatus == model.RequestStatusPending && request.RequestType == model.RequestTypeCreate {
		// Update CREATE request PENDING: cho phép update tất cả fields
		var req dto.RequestCoaAccountCreateRequestWithValidation
		if err := c.ShouldBindJSON(&req); err != nil {
			out := validate.FormatErrorMessage(req, err)
			ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
			return
		}

		accountNo := req.AccountData.AccountNo

		// Check duplicate (loại trừ chính request hiện tại đang được update)
		duplicateInfo, err := h.service.CheckDuplicate(c, accountNo, model.RequestTypeCreate, request.ID)
		if err != nil {
			h.logger.Error("Failed to check duplicate", err)
			ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi kiểm tra trùng lặp")
			return
		}

		if duplicateInfo.IsDuplicate {
			c.JSON(http.StatusBadRequest, dto.PreResponse{
				Data: map[string]interface{}{
					"is_duplicate": true,
					"message":      duplicateInfo.Message,
					"duplicate_in": duplicateInfo.DuplicateIn,
				},
			})
			return
		}

		// Check duplicate code (trừ chính request hiện tại nếu code không đổi)
		currentAccountData, _ := request.GetAccountData()
		if currentAccountData == nil || currentAccountData.Code != req.AccountData.Code {
			exists, err := h.coAccountRepo.Exists(c, map[string]any{
				"code":   req.AccountData.Code,
				"status": "ACTIVE",
			})
			if err != nil {
				ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra")
				return
			}
			if exists {
				ginhp.RespondError(c, http.StatusBadRequest, "Mã tài khoản đã tồn tại")
				return
			}
		}

		// Convert to model và update
		account := &model.CoaAccount{
			Code:        req.AccountData.Code,
			AccountNo:   req.AccountData.AccountNo,
			Name:        req.AccountData.Name,
			Type:        req.AccountData.Type,
			Currency:    req.AccountData.Currency,
			Status:      req.AccountData.Status,
			ParentID:    req.AccountData.ParentID,
			Provider:    req.AccountData.Provider,
			Network:     req.AccountData.Network,
			Description: req.AccountData.Description,
		}

		// Handle Description (store in metadata)
		// if req.AccountData.Description != nil {
		// 	metadata := map[string]interface{}{
		// 		"description": *req.AccountData.Description,
		// 	}
		// 	metadataJSON, err := json.Marshal(metadata)
		// 	if err != nil {
		// 		ginhp.RespondError(c, http.StatusBadRequest, "Failed to marshal description")
		// 		return
		// 	}
		// 	account.Metadata = (*datatypes.JSON)(&metadataJSON)
		// }

		if err := request.SetAccountData(account); err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}

	} else if request.RequestStatus == model.RequestStatusRejected {
		// Update REJECTED request: chỉ được phép edit AccountNo, Status, Description
		var req dto.RequestCoaAccountUpdateRequestWithValidation
		if err := c.ShouldBindJSON(&req); err != nil {
			out := validate.FormatErrorMessage(req, err)
			ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
			return
		}

		// Update account data
		account := &model.CoaAccount{
			AccountNo: req.AccountData.AccountNo,
			Status:    req.AccountData.Status,
		}

		// Handle Description (store in metadata)
		// if req.AccountData.Description != nil {
		// 	metadata := map[string]interface{}{
		// 		"description": *req.AccountData.Description,
		// 	}
		// 	metadataJSON, err := json.Marshal(metadata)
		// 	if err != nil {
		// 		ginhp.RespondError(c, http.StatusBadRequest, "Failed to marshal description")
		// 		return
		// 	}
		// 	account.Metadata = (*datatypes.JSON)(&metadataJSON)
		// }

		if err := request.SetAccountData(account); err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}

		// Reset status về PENDING và clear reject info
		request.RequestStatus = model.RequestStatusPending
		request.CheckerID = nil
		request.RejectReason = nil
		request.Comment = nil
		request.CheckedAt = nil
	} else {
		ginhp.RespondError(c, http.StatusBadRequest, "Chỉ có thể cập nhật request PENDING với type CREATE hoặc request REJECTED")
		return
	}

	if err := h.service.Update(c, request); err != nil {
		h.logger.Error("Failed to update request", err)
		ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi cập nhật request")
		return
	}

	ginhp.RespondOK(c, "Cập nhật request thành công. Đang chờ phê duyệt lại.")
}

// Approve godoc
// @Summary Approve COA account request
// @Description Approve a pending COA account request. This will create or update the COA account in the core system.
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body dto.RequestCoaAccountApproveRequest true "Approve request"
// @Success 200 {object} ginhp.Response
// @Failure 400 {object} dto.PreResponse
// @Failure 401 {object} dto.PreResponse
// @Failure 404 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts/{id}/approve [post]
// Approve handles POST /request-coa-accounts/:id/approve
func (h *RequestCoaAccountHandler) Approve(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Lấy checkerID từ context
	checkerID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.RequestCoaAccountApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
		return
	}

	// Lấy request
	request, err := h.service.GetRequestByID(c, uint64(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusNotFound, "Request not found")
			return
		}
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Approve request (sẽ tự động tạo/cập nhật account trong coa_accounts)
	if err := request.Approve(h.db, uint64(checkerID), req.Comment); err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusBadRequest, "Request không thể được phê duyệt")
			return
		}
		h.logger.Error("Failed to approve request", err)
		ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi phê duyệt request")
		return
	}

	ginhp.RespondOK(c, "Phê duyệt request thành công")
}

// Reject godoc
// @Summary Reject COA account request
// @Description Reject a pending COA account request with a reason
// @Tags request-coa-accounts
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body dto.RequestCoaAccountRejectRequest true "Reject request"
// @Success 200 {object} ginhp.Response
// @Failure 400 {object} dto.PreResponse
// @Failure 401 {object} dto.PreResponse
// @Failure 404 {object} dto.PreResponse
// @Failure 422 {object} dto.PreResponse
// @Failure 500 {object} dto.PreResponse
// @Router /request-coa-accounts/{id}/reject [post]
// Reject handles POST /request-coa-accounts/:id/reject
func (h *RequestCoaAccountHandler) Reject(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Lấy checkerID từ context
	checkerID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.RequestCoaAccountRejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
		return
	}

	// Lấy request
	request, err := h.service.GetRequestByID(c, uint64(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusNotFound, "Request not found")
			return
		}
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Reject request
	if err := request.Reject(h.db, uint64(checkerID), req.RejectReason); err != nil {
		if err == gorm.ErrRecordNotFound {
			ginhp.RespondError(c, http.StatusBadRequest, "Request không thể bị từ chối")
			return
		}
		h.logger.Error("Failed to reject request", err)
		ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi từ chối request")
		return
	}

	ginhp.RespondOK(c, "Từ chối request thành công")
}
