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
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RequestCoaAccountHandler struct {
	logger        logger.CustomLogger
	service       *RequestCoaAccountService
	coAccountRepo repo.CoAccountRepo
	db            *gorm.DB
}

func NewRequestCoaAccountHandler(service *RequestCoaAccountService, db *gorm.DB, coAccountRepo repo.CoAccountRepo) *RequestCoaAccountHandler {
	return &RequestCoaAccountHandler{
		logger:        logger.NewSystemLog("RequestCoaAccountHandler"),
		service:       service,
		coAccountRepo: coAccountRepo,
		db:            db,
	}
}

// GetList handles GET /request-coa-accounts
func (h *RequestCoaAccountHandler) GetList(c *gin.Context) {
	var filter dto.ListRequestCoaAccountFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.GetList(c, &filter)
	if err != nil {
		h.logger.Error("Failed to get request list", err)
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

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
		fmt.Println("AccountId", req.AccountData.AccountId)
		accountNo = req.AccountData.AccountNo
		_, err := h.coAccountRepo.GetOneByFields(c, map[string]interface{}{
			"id": req.AccountData.AccountId,
		})
		if err != nil {
			fmt.Println("AccountId", req.AccountData.AccountId)
			ginhp.ReturnBadRequestError(c, err)
			return
		}

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

// Update handles PUT /request-coa-accounts/:id
// Chỉ cho phép update request có status REJECTED
func (h *RequestCoaAccountHandler) Update(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.RequestCoaAccountUpdateRequestWithValidation
	if err := c.ShouldBindJSON(&req); err != nil {
		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)
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

	// Chỉ cho phép update request bị reject
	if request.RequestStatus != model.RequestStatusRejected {
		ginhp.RespondError(c, http.StatusBadRequest, "Chỉ có thể cập nhật request đã bị từ chối")
		return
	}

	// Update account data
	// For UPDATE: chỉ được phép edit AccountNo, Status, Description (giống EDIT)
	account := &model.CoaAccount{
		AccountNo: req.AccountData.AccountNo,
		Status:    req.AccountData.Status,
	}

	// Handle Description (store in metadata)
	if req.AccountData.Description != nil {
		metadata := map[string]interface{}{
			"description": *req.AccountData.Description,
		}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			ginhp.RespondError(c, http.StatusBadRequest, "Failed to marshal description")
			return
		}
		account.Metadata = (*datatypes.JSON)(&metadataJSON)
	}

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

	if err := h.service.Update(c, request); err != nil {
		h.logger.Error("Failed to update request", err)
		ginhp.RespondError(c, http.StatusInternalServerError, "Có lỗi xảy ra khi cập nhật request")
		return
	}

	ginhp.RespondOK(c, "Cập nhật request thành công. Đang chờ phê duyệt lại.")
}

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
