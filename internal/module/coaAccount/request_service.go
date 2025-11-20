package coaaccount

import (
	"context"
	"core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type RequestCoaAccountService struct {
	logger        logger.CustomLogger
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
}

type DuplicateInfo struct {
	IsDuplicate bool   `json:"is_duplicate"`
	Message     string `json:"message"`
	DuplicateIn string `json:"duplicate_in"` // "request" hoặc "core"
}

func NewRequestCoaAccountService(db *gorm.DB, coAccountRepo repo.CoAccountRepo) *RequestCoaAccountService {
	return &RequestCoaAccountService{
		logger:        logger.NewSystemLog("RequestCoaAccountService"),
		db:            db,
		coAccountRepo: coAccountRepo,
	}
}

// CheckDuplicate checks if account_no already exists in request_coa_accounts or coa_accounts
// Theo luồng: Check cả "Đã có record submit" và "Đã có account trong Core"
func (s *RequestCoaAccountService) CheckDuplicate(ctx context.Context, accountNo string, requestType model.RequestType) (*DuplicateInfo, error) {
	info := &DuplicateInfo{
		IsDuplicate: false,
	}

	// 1. Check trong coa_accounts (Core)
	exists, err := s.coAccountRepo.Exists(ctx, map[string]interface{}{
		"account_no": accountNo,
	})
	if err != nil {
		return nil, err
	}
	if exists {
		info.IsDuplicate = true
		info.Message = fmt.Sprintf("Account No. %s đã tồn tại trong Core", accountNo)
		info.DuplicateIn = "core"
		return info, nil
	}

	// 2. Check trong request_coa_accounts với status PENDING hoặc APPROVED
	// Sử dụng JSONB query để check account_no trong data field
	// PostgreSQL JSONB: data->>'account_no' để lấy giá trị string
	var count int64
	query := s.db.WithContext(ctx).Model(&model.RequestCoaAccount{}).
		Where("data->>'account_no' = ?", accountNo).
		Where("request_status IN ?", []string{
			string(model.RequestStatusPending),
			string(model.RequestStatusApproved),
		})

	if err := query.Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		info.IsDuplicate = true
		info.Message = fmt.Sprintf("Account No. %s đã có record đang chờ phê duyệt hoặc đã được phê duyệt", accountNo)
		info.DuplicateIn = "request"
		return info, nil
	}

	return info, nil
}

// GetList returns paginated list of requests
func (s *RequestCoaAccountService) GetList(ctx context.Context, filter *dto.ListRequestCoaAccountFilter) (*dto.PaginationResponse[*dto.RequestCoaAccountResponse], error) {
	query := s.db.WithContext(ctx).Model(&model.RequestCoaAccount{})

	// Apply filters
	if filter.RequestType != nil {
		query = query.Where("request_type = ?", *filter.RequestType)
	}
	if filter.RequestStatus != nil {
		query = query.Where("request_status = ?", *filter.RequestStatus)
	}
	if filter.MakerID != nil {
		query = query.Where("maker_id = ?", *filter.MakerID)
	}
	if filter.CheckerID != nil {
		query = query.Where("checker_id = ?", *filter.CheckerID)
	}
	if filter.CoaAccountID != nil {
		query = query.Where("coa_account_id = ?", *filter.CoaAccountID)
	}
	if filter.Search != nil && *filter.Search != "" {
		search := "%" + *filter.Search + "%"
		query = query.Where(
			"data->>'account_no' ILIKE ? OR data->>'name' ILIKE ? OR data->>'code' ILIKE ?",
			search, search, search,
		)
	}

	// Preload relations
	query = query.Preload("CoaAccount").Preload("Maker").Preload("Checker")

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	page := int64(1)
	if filter.Page != nil && *filter.Page > 0 {
		page = *filter.Page
	}
	limit := int64(10)
	if filter.Limit != nil && *filter.Limit > 0 {
		limit = *filter.Limit
	}
	offset := (page - 1) * limit

	// Get data
	var requests []model.RequestCoaAccount
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(limit)).Find(&requests).Error; err != nil {
		return nil, err
	}

	// Convert to response DTOs
	items := make([]*dto.RequestCoaAccountResponse, 0, len(requests))
	for _, req := range requests {
		item, err := s.toResponseDTO(&req)
		if err != nil {
			s.logger.Error("Failed to convert request to DTO", err)
			continue
		}
		items = append(items, item)
	}

	totalPage := (total + limit - 1) / limit
	var nextPage *int64
	var prevPage *int64
	if page < totalPage {
		np := page + 1
		nextPage = &np
	}
	if page > 1 {
		pp := page - 1
		prevPage = &pp
	}

	return &dto.PaginationResponse[*dto.RequestCoaAccountResponse]{
		Items:     items,
		Total:     total,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}

// GetByID returns a request by ID
func (s *RequestCoaAccountService) GetByID(ctx context.Context, id uint64) (*dto.RequestCoaAccountResponse, error) {
	var request model.RequestCoaAccount
	if err := s.db.WithContext(ctx).
		Preload("CoaAccount").
		Preload("Maker").
		Preload("Checker").
		First(&request, id).Error; err != nil {
		return nil, err
	}

	return s.toResponseDTO(&request)
}

// GetRequestByID returns the model (not DTO)
func (s *RequestCoaAccountService) GetRequestByID(ctx context.Context, id uint64) (*model.RequestCoaAccount, error) {
	var request model.RequestCoaAccount
	if err := s.db.WithContext(ctx).First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// Create creates a new request
func (s *RequestCoaAccountService) Create(ctx context.Context, request *model.RequestCoaAccount) error {
	return s.db.WithContext(ctx).Create(request).Error
}

// Update updates an existing request
func (s *RequestCoaAccountService) Update(ctx context.Context, request *model.RequestCoaAccount) error {
	return s.db.WithContext(ctx).Save(request).Error
}

// toResponseDTO converts model to response DTO
func (s *RequestCoaAccountService) toResponseDTO(req *model.RequestCoaAccount) (*dto.RequestCoaAccountResponse, error) {
	// Parse data from JSON
	var dataMap map[string]interface{}
	dataBytes := []byte(req.Data)
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("failed to parse request data: %w", err)
	}

	response := &dto.RequestCoaAccountResponse{
		ID:            req.ID,
		CoaAccountID:  req.CoaAccountID,
		RequestType:   string(req.RequestType),
		RequestStatus: string(req.RequestStatus),
		MakerID:       req.MakerID,
		CheckerID:     req.CheckerID,
		Data:          dataMap,
		Comment:       req.Comment,
		RejectReason:  req.RejectReason,
		CreatedAt:     req.CreatedAt,
		UpdatedAt:     req.UpdatedAt,
		CheckedAt:     req.CheckedAt,
		CoaAccount:    req.CoaAccount,
		Maker:         req.Maker,
		Checker:       req.Checker,
	}

	return response, nil
}

