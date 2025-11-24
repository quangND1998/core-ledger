package dto

import (
	model "core-ledger/model/core-ledger"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

// RequestCoaAccountCreateRequest DTO for creating a request
type RequestCoaAccountCreateRequest struct {
	RequestType model.RequestType `json:"request_type" binding:"required,oneof=CREATE EDIT"`
	AccountData *CoaAccountData   `json:"account_data" binding:"required"`
}

// RequestCoaAccountCreateRequestWithValidation DTO for creating CREATE request with validation
type RequestCoaAccountCreateRequestWithValidation struct {
	RequestType model.RequestType    `json:"request_type" binding:"required,eq=CREATE"`
	AccountData CoaAccountDataCreate `json:"account_data" binding:"required"`
}

// RequestCoaAccountEditRequestWithValidation DTO for creating EDIT request with validation
type RequestCoaAccountEditRequestWithValidation struct {
	RequestType model.RequestType  `json:"request_type" binding:"required,eq=EDIT"`
	AccountData CoaAccountDataEdit `json:"account_data" binding:"required"`
}

// RequestCoaAccountUpdateRequest DTO for updating a rejected request
type RequestCoaAccountUpdateRequest struct {
	AccountData *CoaAccountData `json:"account_data" binding:"required"`
}

// RequestCoaAccountUpdateRequestWithValidation DTO for updating a rejected request with validation
type RequestCoaAccountUpdateRequestWithValidation struct {
	AccountData CoaAccountDataEdit `json:"account_data" binding:"required"`
}

// RequestCoaAccountApproveRequest DTO for approving a request
type RequestCoaAccountApproveRequest struct {
	Comment *string `json:"comment,omitempty"`
}

// RequestCoaAccountRejectRequest DTO for rejecting a request
type RequestCoaAccountRejectRequest struct {
	RejectReason string `json:"reject_reason" binding:"required"`
}

// CoaAccountData represents the account data in request
// For CREATE: all fields are required
// For EDIT: only AccountNo, Status, and Description are allowed to be changed
type CoaAccountData struct {
	// Fields for CREATE (required)
	Code      string  `json:"code,omitempty" ` // Required for CREATE
	AccountNo string  `json:"account_no" binding:"required"`
	Name      string  `json:"name,omitempty"`     // Required for CREATE
	Type      string  `json:"type,omitempty"`     // Required for CREATE
	Currency  string  `json:"currency,omitempty"` // Required for CREATE
	ParentID  *uint64 `json:"parent_id,omitempty"`
	Provider  *string `json:"provider,omitempty"`
	Network   *string `json:"network,omitempty"`

	// Fields for EDIT (allowed to change)
	Status      string  `json:"status,omitempty"`      // Allowed for EDIT
	Description *string `json:"description,omitempty"` // Allowed for EDIT (stored in metadata)
}

// CoaAccountDataCreate represents the account data for CREATE request with validation
type CoaAccountDataCreate struct {
	// Fields for CREATE (required)
	Code      string  `json:"code" binding:"required"`                                 // Required for CREATE
	AccountNo string  `json:"account_no" binding:"required"`                           // Required for CREATE
	Name      string  `json:"name" binding:"required"`                                 // Required for CREATE
	Type      string  `json:"type" binding:"required,oneof=ASSET LIAB EQUITY REV EXP"` // Required for CREATE, must be one of: ASSET, LIAB, EQUITY, REV, EXP
	Currency  string  `json:"currency" binding:"required"`                             // Required for CREATE
	ParentID  *uint64 `json:"parent_id,omitempty"`
	Provider  *string `json:"provider,omitempty"`
	Network   *string `json:"network,omitempty"`

	// Fields for EDIT (allowed to change)
	Status      string  `json:"status" binding:"required"` // Required for CREATE
	Description *string `json:"description,omitempty"`     // Allowed for CREATE (stored in metadata)
}

// CoaAccountDataEdit represents the account data for EDIT request with validation
type CoaAccountDataEdit struct {
	// Fields for EDIT (allowed to change)
	Name        string  `json:"name" binding:"required"`
	AccountId   uint64  `json:"account_id" binding:"required"` // Required for EDIT
	AccountNo   string  `json:"account_no" binding:"required"` // Required for EDIT
	Status      string  `json:"status" binding:"required"`     // Required for EDIT
	Description *string `json:"description,omitempty"`         // Allowed for EDIT (stored in metadata)
}

// ListRequestCoaAccountFilter filter for listing requests
type ListRequestCoaAccountFilter struct {
	BasePaginationQuery
	RequestType   *string `form:"request_type"`   // CREATE, EDIT
	RequestStatus *string `form:"request_status"` // PENDING, APPROVED, REJECTED
	MakerID       *uint64 `form:"maker_id"`
	CheckerID     *uint64 `form:"checker_id"`
	CoaAccountID  *uint64 `form:"coa_account_id"`
	Search        *string `form:"search"` // Search by account_no, name, code
}

// CodeAnalysis phân tích code từ rules - format để frontend tự generate form
type CodeAnalysis struct {
	Code      string `json:"code"`
	TypeCode  string `json:"type_code"`
	TypeName  string `json:"type_name,omitempty"`
	GroupCode string `json:"group_code"`
	GroupName string `json:"group_name,omitempty"`
	// Rules structure với giá trị đã chọn - frontend có thể dùng để render form
	FormData *CodeFormData `json:"form_data,omitempty"`
	IsValid  bool          `json:"is_valid"`
	Error    string        `json:"error,omitempty"`
}

// CodeFormData cấu trúc form data để frontend render
// Format đơn giản: type -> group (đã chọn) -> steps với current_value
type CodeFormData struct {
	Type  CodeFormType  `json:"type"`  // Type với group đã được chọn
	Group CodeFormGroup `json:"group"` // Group đã được chọn với steps đầy đủ
}

// CodeFormType type với group đã chọn
type CodeFormType struct {
	ID        uint64        `json:"id"`
	Code      string        `json:"code"`
	Name      string        `json:"name"`
	Separator string        `json:"separator"`
	Group     CodeFormGroup `json:"group"` // Group đã được chọn
}

// CodeFormGroup group với giá trị đã chọn
type CodeFormGroup struct {
	ID           uint64         `json:"id"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	InputType    string         `json:"input_type"`
	Separator    string         `json:"separator"`
	Selected     bool           `json:"selected"`                // Đánh dấu group này đã được chọn
	CurrentValue string         `json:"current_value,omitempty"` // Giá trị text input của group (nếu input_type là TEXT)
	Steps        []CodeFormStep `json:"steps"`
}

// CodeFormStep step với giá trị đã chọn và options
type CodeFormStep struct {
	StepID           uint64          `json:"step_id"`
	StepOrder        int             `json:"step_order"`
	Type             string          `json:"type"` // SELECT hoặc TEXT
	Label            string          `json:"label,omitempty"`
	CategoryCode     string          `json:"category_code,omitempty"`
	InputCode        string          `json:"input_code,omitempty"`
	InputType        string          `json:"input_type,omitempty"`
	Separator        string          `json:"separator"`
	Values           []RuleValueResp `json:"values,omitempty"`             // Options cho SELECT
	CurrentValue     string          `json:"current_value,omitempty"`      // Giá trị hiện tại đã chọn
	CurrentValueName string          `json:"current_value_name,omitempty"` // Tên của giá trị hiện tại
}

// RequestCoaAccountResponse response DTO
type RequestCoaAccountResponse struct {
	ID            uint64                 `json:"id"`
	CoaAccountID  *uint64                `json:"coa_account_id,omitempty"`
	RequestType   string                 `json:"request_type"`
	RequestStatus string                 `json:"request_status"`
	MakerID       uint64                 `json:"maker_id"`
	CheckerID     *uint64                `json:"checker_id,omitempty"`
	Data          map[string]interface{} `json:"data"`
	Comment       *string                `json:"comment,omitempty"`
	RejectReason  *string                `json:"reject_reason,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CheckedAt     *time.Time             `json:"checked_at,omitempty"`

	// Relations
	CoaAccount *model.CoaAccount `json:"coa_account,omitempty"`
	Maker      *model.User       `json:"maker,omitempty"`
	Checker    *model.User       `json:"checker,omitempty"`

	// Code analysis
	CodeAnalysis *CodeAnalysis `json:"code_analysis,omitempty"`
}

// ToModel converts DTO to model
func (r *RequestCoaAccountCreateRequest) ToModel(makerID uint64) (*model.RequestCoaAccount, error) {
	account := &model.CoaAccount{
		AccountNo: r.AccountData.AccountNo,
		Status:    r.AccountData.Status,
	}

	// For CREATE: set all required fields
	if r.RequestType == model.RequestTypeCreate {
		if r.AccountData.Code == "" || r.AccountData.Name == "" || r.AccountData.Type == "" || r.AccountData.Currency == "" {
			return nil, fmt.Errorf("code, name, type, and currency are required for CREATE request")
		}
		account.Code = r.AccountData.Code
		account.Name = r.AccountData.Name
		account.Type = r.AccountData.Type
		account.Currency = r.AccountData.Currency
		account.ParentID = r.AccountData.ParentID
		account.Provider = r.AccountData.Provider
		account.Network = r.AccountData.Network
	}

	// Handle Description (store in metadata)
	if r.AccountData.Description != nil {
		metadata := map[string]interface{}{
			"description": *r.AccountData.Description,
		}
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal description: %w", err)
		}
		account.Metadata = (*datatypes.JSON)(&metadataJSON)
	}

	var request *model.RequestCoaAccount
	var err error

	if r.RequestType == model.RequestTypeCreate {
		request, err = model.NewCreateRequest(account, makerID)
	} else {
		request, err = model.NewEditRequest(account, makerID)
	}

	if err != nil {
		return nil, err
	}

	return request, nil
}

// ToModel converts DTO to model for CREATE request with validation
func (r *RequestCoaAccountCreateRequestWithValidation) ToModel(makerID uint64) (*model.RequestCoaAccount, error) {
	account := &model.CoaAccount{
		Code:        r.AccountData.Code,
		AccountNo:   r.AccountData.AccountNo,
		Name:        r.AccountData.Name,
		Type:        r.AccountData.Type,
		Currency:    r.AccountData.Currency,
		Status:      r.AccountData.Status,
		ParentID:    r.AccountData.ParentID,
		Provider:    r.AccountData.Provider,
		Network:     r.AccountData.Network,
		Description: r.AccountData.Description,
	}

	// Handle Description (store in metadata)
	// if r.AccountData.Description != nil {
	// 	metadata := map[string]interface{}{
	// 		"description": *r.AccountData.Description,
	// 	}
	// 	metadataJSON, err := json.Marshal(metadata)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to marshal description: %w", err)
	// 	}
	// 	account.Metadata = (*datatypes.JSON)(&metadataJSON)
	// }

	request, err := model.NewCreateRequest(account, makerID)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ToModel converts DTO to model for EDIT request with validation
func (r *RequestCoaAccountEditRequestWithValidation) ToModel(makerID uint64) (*model.RequestCoaAccount, error) {
	account := &model.CoaAccount{
		ID:          r.AccountData.AccountId,
		Name:        r.AccountData.Name,
		Description: r.AccountData.Description,
		AccountNo:   r.AccountData.AccountNo,
		Status:      r.AccountData.Status,
	}

	// Handle Description (store in metadata)
	// if r.AccountData.Description != nil {
	// 	metadata := map[string]interface{}{
	// 		"description": *r.AccountData.Description,
	// 	}
	// 	metadataJSON, err := json.Marshal(metadata)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to marshal description: %w", err)
	// 	}
	// 	account.Metadata = (*datatypes.JSON)(&metadataJSON)
	// }

	request, err := model.NewEditRequest(account, makerID)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ToModel converts DTO to model for UPDATE request with validation
func (r *RequestCoaAccountUpdateRequestWithValidation) ToModel() (*CoaAccountData, error) {
	accountData := &CoaAccountData{
		AccountNo:   r.AccountData.AccountNo,
		Status:      r.AccountData.Status,
		Description: r.AccountData.Description,
	}
	return accountData, nil
}
