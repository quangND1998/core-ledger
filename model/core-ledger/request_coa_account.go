package model

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RequestType represents the type of request
type RequestType string

const (
	RequestTypeCreate RequestType = "CREATE"
	RequestTypeEdit   RequestType = "EDIT"
)

// RequestStatus represents the status of request
type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "PENDING"
	RequestStatusApproved RequestStatus = "APPROVED"
	RequestStatusRejected RequestStatus = "REJECTED"
)

// RequestCoaAccount represents a request to create or edit a COA account
type RequestCoaAccount struct {
	Entity
	ID          uint64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CoaAccountID *uint64        `gorm:"column:coa_account_id;index" json:"coa_account_id,omitempty"` // NULL nếu là CREATE
	RequestType  RequestType    `gorm:"type:varchar(16);not null;check:request_type IN ('CREATE','EDIT')" json:"request_type"`
	RequestStatus RequestStatus `gorm:"type:varchar(16);not null;default:'PENDING';check:request_status IN ('PENDING','APPROVED','REJECTED')" json:"request_status"`
	
	// User tracking
	MakerID   uint64  `gorm:"column:maker_id;not null;index" json:"maker_id"`   // User tạo request
	CheckerID *uint64 `gorm:"column:checker_id;index" json:"checker_id,omitempty"` // User approve/reject
	
	// Data thay đổi (lưu toàn bộ data của CoaAccount)
	Data datatypes.JSON `gorm:"type:jsonb;not null" json:"data"` // Lưu toàn bộ fields của CoaAccount
	
	// Optional fields
	Comment *string `gorm:"type:text" json:"comment,omitempty"` // Comment từ checker
	RejectReason *string `gorm:"type:text" json:"reject_reason,omitempty"` // Lý do reject
	
	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CheckedAt *time.Time `gorm:"column:checked_at" json:"checked_at,omitempty"` // Thời gian approve/reject

	// Relations
	CoaAccount *CoaAccount `gorm:"foreignKey:CoaAccountID" json:"coa_account,omitempty"`
	Maker      *User       `gorm:"foreignKey:MakerID" json:"maker,omitempty"`
	Checker    *User       `gorm:"foreignKey:CheckerID" json:"checker,omitempty"`
}

func (r *RequestCoaAccount) TableName() string {
	return "request_coa_accounts"
}

// BeforeCreate hook
func (r *RequestCoaAccount) BeforeCreate(tx *gorm.DB) error {
	if r.RequestStatus == "" {
		r.RequestStatus = RequestStatusPending
	}
	return nil
}

// IsPending checks if request is pending
func (r *RequestCoaAccount) IsPending() bool {
	return r.RequestStatus == RequestStatusPending
}

// IsApproved checks if request is approved
func (r *RequestCoaAccount) IsApproved() bool {
	return r.RequestStatus == RequestStatusApproved
}

// IsRejected checks if request is rejected
func (r *RequestCoaAccount) IsRejected() bool {
	return r.RequestStatus == RequestStatusRejected
}

// CanApprove checks if request can be approved
func (r *RequestCoaAccount) CanApprove() bool {
	return r.IsPending()
}

// CanReject checks if request can be rejected
func (r *RequestCoaAccount) CanReject() bool {
	return r.IsPending()
}

// Approve approves the request and applies changes to coa_accounts
func (r *RequestCoaAccount) Approve(db *gorm.DB, checkerID uint64, comment *string) error {
	if !r.CanApprove() {
		return gorm.ErrRecordNotFound
	}

	// Parse data từ JSON
	accountData, err := r.GetAccountData()
	if err != nil {
		return err
	}

	// Bắt đầu transaction
	return db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		if r.RequestType == RequestTypeCreate {
			// Tạo mới COA account
			accountData.ID = 0 // Reset ID để tạo mới
			if err := tx.Create(accountData).Error; err != nil {
				return err
			}
			// Cập nhật coa_account_id sau khi tạo
			r.CoaAccountID = &accountData.ID
		} else if r.RequestType == RequestTypeEdit {
			// Cập nhật COA account hiện có
			if r.CoaAccountID == nil {
				return gorm.ErrRecordNotFound
			}
			
			// Lấy account hiện tại
			var existingAccount CoaAccount
			if err := tx.First(&existingAccount, *r.CoaAccountID).Error; err != nil {
				return err
			}

			// Khi Edit, chỉ được phép update: AccountNo, Status, Description
			// Giữ nguyên các field khác (Code, Name, Type, Currency, etc.)
			existingAccount.AccountNo = accountData.AccountNo
			existingAccount.Status = accountData.Status
			existingAccount.UpdatedAt = now

			// Update Description trong metadata nếu có
			// Nếu accountData có Metadata, merge với existing metadata
			if accountData.Metadata != nil {
				// Parse existing metadata
				var existingMetadata map[string]interface{}
				if existingAccount.Metadata != nil {
					if err := json.Unmarshal([]byte(*existingAccount.Metadata), &existingMetadata); err != nil {
						existingMetadata = make(map[string]interface{})
					}
				} else {
					existingMetadata = make(map[string]interface{})
				}

				// Parse new metadata
				var newMetadata map[string]interface{}
				if err := json.Unmarshal([]byte(*accountData.Metadata), &newMetadata); err == nil {
					// Merge description từ newMetadata vào existingMetadata
					if desc, ok := newMetadata["description"]; ok {
						existingMetadata["description"] = desc
					}
					// Marshal lại
					mergedJSON, err := json.Marshal(existingMetadata)
					if err == nil {
						merged := datatypes.JSON(mergedJSON)
						existingAccount.Metadata = &merged
					}
				}
			}

			if err := tx.Save(&existingAccount).Error; err != nil {
				return err
			}
		}

		// Cập nhật request status
		r.RequestStatus = RequestStatusApproved
		r.CheckerID = &checkerID
		r.Comment = comment
		r.CheckedAt = &now

		if err := tx.Save(r).Error; err != nil {
			return err
		}

		return nil
	})
}

// Reject rejects the request
func (r *RequestCoaAccount) Reject(db *gorm.DB, checkerID uint64, rejectReason string) error {
	if !r.CanReject() {
		return gorm.ErrRecordNotFound
	}

	now := time.Now()
	r.RequestStatus = RequestStatusRejected
	r.CheckerID = &checkerID
	r.RejectReason = &rejectReason
	r.CheckedAt = &now

	return db.Save(r).Error
}

// GetAccountData parses and returns the CoaAccount data from JSON
func (r *RequestCoaAccount) GetAccountData() (*CoaAccount, error) {
	var accountData CoaAccount
	dataBytes := []byte(r.Data)
	if err := json.Unmarshal(dataBytes, &accountData); err != nil {
		return nil, err
	}
	return &accountData, nil
}

// SetAccountData sets the CoaAccount data as JSON
func (r *RequestCoaAccount) SetAccountData(account *CoaAccount) error {
	dataBytes, err := json.Marshal(account)
	if err != nil {
		return err
	}
	r.Data = datatypes.JSON(dataBytes)
	return nil
}

// NewCreateRequest creates a new CREATE request for a COA account
func NewCreateRequest(account *CoaAccount, makerID uint64) (*RequestCoaAccount, error) {
	request := &RequestCoaAccount{
		RequestType:   RequestTypeCreate,
		RequestStatus: RequestStatusPending,
		MakerID:       makerID,
		CoaAccountID:  nil, // NULL for CREATE
	}
	
	if err := request.SetAccountData(account); err != nil {
		return nil, err
	}
	
	return request, nil
}

// NewEditRequest creates a new EDIT request for an existing COA account
func NewEditRequest(account *CoaAccount, makerID uint64) (*RequestCoaAccount, error) {
	if account.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	
	request := &RequestCoaAccount{
		RequestType:   RequestTypeEdit,
		RequestStatus: RequestStatusPending,
		MakerID:       makerID,
		CoaAccountID:  &account.ID,
	}
	
	if err := request.SetAccountData(account); err != nil {
		return nil, err
	}
	
	return request, nil
}

