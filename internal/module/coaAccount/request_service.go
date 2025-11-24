package coaaccount

import (
	"context"
	"core-ledger/internal/module/ruleCategory"
	"core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequestCoaAccountService struct {
	logger           logger.CustomLogger
	db               *gorm.DB
	coAccountRepo    repo.CoAccountRepo
	ruleCategoryService *ruleCategory.RuleCateogySerive
}

type DuplicateInfo struct {
	IsDuplicate bool   `json:"is_duplicate"`
	Message     string `json:"message"`
	DuplicateIn string `json:"duplicate_in"` // "request" hoặc "core"
}

func NewRequestCoaAccountService(db *gorm.DB, coAccountRepo repo.CoAccountRepo, ruleCategoryService *ruleCategory.RuleCateogySerive) *RequestCoaAccountService {
	return &RequestCoaAccountService{
		logger:              logger.NewSystemLog("RequestCoaAccountService"),
		db:                  db,
		coAccountRepo:       coAccountRepo,
		ruleCategoryService: ruleCategoryService,
	}
}

// CheckDuplicate checks if account_no already exists in request_coa_accounts or coa_accounts
// Theo luồng: Check cả "Đã có record submit" và "Đã có account trong Core"
// excludeRequestID: ID của request cần loại trừ khỏi kiểm tra (ví dụ: khi update request hiện tại)
func (s *RequestCoaAccountService) CheckDuplicate(ctx context.Context, accountNo string, requestType model.RequestType, excludeRequestID ...uint64) (*DuplicateInfo, error) {
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

	// Loại trừ request hiện tại nếu có excludeRequestID
	if len(excludeRequestID) > 0 && excludeRequestID[0] > 0 {
		query = query.Where("id != ?", excludeRequestID[0])
	}

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

	// Convert to response DTOs (không parse code_analysis để tối ưu performance)
	items := make([]*dto.RequestCoaAccountResponse, 0, len(requests))
	for _, req := range requests {
		item, err := s.toResponseDTO(ctx, &req, false) // false = không parse code_analysis
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

	return s.toResponseDTO(ctx, &request, true) // true = có parse code_analysis
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

// parseCode phân tích code dựa trên rules
func (s *RequestCoaAccountService) parseCode(ctx context.Context, code string) (*dto.CodeAnalysis, error) {
	if code == "" {
		return nil, nil
	}

	// Kiểm tra ruleCategoryService
	if s.ruleCategoryService == nil {
		return nil, fmt.Errorf("ruleCategoryService is nil")
	}

	// Tạo gin context từ context để gọi GetCoaAccountRules
	ginCtx, _ := gin.CreateTestContext(nil)
	// Tạo request mới với context
	req, _ := http.NewRequestWithContext(ctx, "GET", "/", nil)
	ginCtx.Request = req

	// Lấy rules từ service
	rules, err := s.ruleCategoryService.GetCoaAccountRules(ginCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}

	// Tìm type code (phần đầu trước separator đầu tiên)
	var matchedType *dto.CoaAccountRuleTypeResp
	typeCode := ""
	typeSeparator := ""

	for i := range rules {
		ruleType := rules[i]
		// Thử tách code theo separator của type
		if ruleType.Separator != "" && strings.HasPrefix(code, ruleType.Code+ruleType.Separator) {
			typeCode = ruleType.Code
			typeSeparator = ruleType.Separator
			matchedType = &rules[i]
			break
		}
	}

	if matchedType == nil {
		return &dto.CodeAnalysis{
			Code:      code,
			TypeCode:  "",
			GroupCode: "",
			FormData:  nil,
			IsValid:   false,
			Error:     "Không tìm thấy type code phù hợp",
		}, nil
	}

	// Tách phần còn lại sau type code
	remainingCode := strings.TrimPrefix(code, typeCode+typeSeparator)

	// Tìm group code
	var matchedGroup *dto.CoaAccountRuleGroupResp
	groupCode := ""
	groupSeparator := ""

	// Trường hợp 1: Type có groups (ASSET, LIAB, etc.)
	if len(matchedType.Groups) > 0 {
		// Tìm group trong code (có group code prefix)
		for i := range matchedType.Groups {
			group := matchedType.Groups[i]
			if group.Separator != "" && strings.HasPrefix(remainingCode, group.Code+group.Separator) {
				groupCode = group.Code
				groupSeparator = group.Separator
				matchedGroup = &matchedType.Groups[i]
				break
			}
		}

		// Trường hợp 2: Không tìm thấy group code trong string
		// Có thể là: REV/EXP với group KIND, hoặc LIAB với group DETAILS (không có group code trong string)
		if matchedGroup == nil && remainingCode != "" {
			for i := range matchedType.Groups {
				group := matchedType.Groups[i]
				
				// Nếu là group KIND (REV/EXP) - có thể có hoặc không có group code prefix
				if group.Code == "KIND" {
					if group.Separator != "" && strings.HasPrefix(remainingCode, group.Code+group.Separator) {
						// Có group code prefix: "KIND:..."
						groupCode = group.Code
						groupSeparator = group.Separator
						matchedGroup = &matchedType.Groups[i]
						break
					} else if len(group.Steps) > 0 {
						// Không có group code prefix, parse trực tiếp: "REV:SERVICE_FEE.USD"
						groupCode = group.Code
						groupSeparator = ""
						matchedGroup = &matchedType.Groups[i]
						break
					}
				}
				
				// Nếu type chỉ có 1 group duy nhất và không tìm thấy group code trong string
				// Ví dụ: LIAB với group DETAILS (input_type TEXT) - "LIAB:123123:USD"
				if len(matchedType.Groups) == 1 && len(group.Steps) > 0 {
					groupCode = group.Code
					groupSeparator = "" // Không có group code trong string
					matchedGroup = &matchedType.Groups[i]
					break
				}
			}
		}
	}

	if matchedGroup == nil {
		return &dto.CodeAnalysis{
			Code:      code,
			TypeCode:  typeCode,
			TypeName:  matchedType.Name,
			GroupCode: "",
			FormData:  nil,
			IsValid:   false,
			Error:     "Không tìm thấy group code phù hợp",
		}, nil
	}

	// Tách phần còn lại sau group code (nếu có group code trong string)
	if groupSeparator != "" && strings.HasPrefix(remainingCode, groupCode+groupSeparator) {
		remainingCode = strings.TrimPrefix(remainingCode, groupCode+groupSeparator)
	}

	// Parse group text input (nếu group có input_type TEXT)
	groupTextValue := ""
	currentCode := remainingCode

	// Nếu group có input_type TEXT, parse text input của group trước
	if matchedGroup.InputType == "TEXT" && matchedGroup.Separator != "" {
		// Tách text input của group theo group separator
		if strings.Contains(currentCode, matchedGroup.Separator) {
			parts := strings.SplitN(currentCode, matchedGroup.Separator, 2)
			groupTextValue = parts[0]
			if len(parts) > 1 {
				currentCode = parts[1]
			} else {
				currentCode = ""
			}
		} else {
			// Nếu không có separator, toàn bộ là text input của group
			// Nhưng cần kiểm tra xem có steps không
			if len(matchedGroup.Steps) > 0 {
				// Có steps, cần parse tiếp
				// Text input của group là phần trước step đầu tiên
				firstStep := matchedGroup.Steps[0]
				if firstStep.Separator != "" && strings.Contains(currentCode, firstStep.Separator) {
					parts := strings.SplitN(currentCode, firstStep.Separator, 2)
					groupTextValue = parts[0]
					currentCode = parts[1]
				} else {
					// Không có separator của step, toàn bộ là text input
					groupTextValue = currentCode
					currentCode = ""
				}
			} else {
				// Không có steps, toàn bộ là text input
				groupTextValue = currentCode
				currentCode = ""
			}
		}
	}

	// Parse các steps và lưu giá trị đã chọn
	stepValues := make(map[int]string) // step_order -> value

	// Tách các step theo separator của từng step
	for _, step := range matchedGroup.Steps {
		if currentCode == "" {
			break
		}

		stepValue := ""
		stepSeparator := step.Separator

		if stepSeparator != "" && strings.Contains(currentCode, stepSeparator) {
			// Tách theo separator
			parts := strings.SplitN(currentCode, stepSeparator, 2)
			stepValue = parts[0]
			if len(parts) > 1 {
				currentCode = parts[1]
			} else {
				currentCode = ""
			}
		} else {
			// Step cuối cùng hoặc không có separator
			stepValue = currentCode
			currentCode = ""
		}

		stepValues[step.StepOrder] = stepValue
	}

	// Build form data với cấu trúc rules và giá trị đã chọn
	formSteps := []dto.CodeFormStep{}
	for _, step := range matchedGroup.Steps {
		currentValue := stepValues[step.StepOrder]
		currentValueName := ""

		// Nếu là SELECT, tìm value name từ values
		if step.Type == "SELECT" && len(step.Values) > 0 && currentValue != "" {
			for _, val := range step.Values {
				if val.Value == currentValue {
					currentValueName = val.Name
					break
				}
			}
		}

		formStep := dto.CodeFormStep{
			StepID:           step.StepID,
			StepOrder:        step.StepOrder,
			Type:             step.Type,
			Label:            step.Label,
			CategoryCode:     step.CategoryCode,
			InputCode:        step.InputCode,
			InputType:        step.InputType,
			Separator:        step.Separator,
			Values:           step.Values, // Tất cả options
			CurrentValue:     currentValue,
			CurrentValueName: currentValueName,
		}

		formSteps = append(formSteps, formStep)
	}

	// Build selected group với đầy đủ steps và current values
	selectedGroup := dto.CodeFormGroup{
		ID:           matchedGroup.ID,
		Code:         matchedGroup.Code,
		Name:         matchedGroup.Name,
		InputType:    matchedGroup.InputType,
		Separator:    matchedGroup.Separator,
		CurrentValue: groupTextValue, // Text input của group (nếu input_type là TEXT)
		Steps:        formSteps,      // Steps với current values đã được parse
	}

	// Build type với group đã chọn
	formType := dto.CodeFormType{
		ID:        matchedType.ID,
		Code:      matchedType.Code,
		Name:      matchedType.Name,
		Separator: matchedType.Separator,
		Group:     selectedGroup,
	}

	formData := &dto.CodeFormData{
		Type:  formType,
		Group: selectedGroup,
	}

	return &dto.CodeAnalysis{
		Code:      code,
		TypeCode:  typeCode,
		TypeName:  matchedType.Name,
		GroupCode: groupCode,
		GroupName: matchedGroup.Name,
		FormData:  formData,
		IsValid:   true,
	}, nil
}

// toResponseDTO converts model to response DTO
// includeCodeAnalysis: nếu true thì sẽ parse code_analysis (dùng cho GetDetail), false thì bỏ qua (dùng cho GetList để tối ưu)
func (s *RequestCoaAccountService) toResponseDTO(ctx context.Context, req *model.RequestCoaAccount, includeCodeAnalysis bool) (*dto.RequestCoaAccountResponse, error) {
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

	// Chỉ phân tích code nếu includeCodeAnalysis = true (để tối ưu performance cho GetList)
	if includeCodeAnalysis {
		if code, ok := dataMap["code"].(string); ok && code != "" {
			codeAnalysis, err := s.parseCode(ctx, code)
			if err != nil {
				s.logger.Error("Failed to parse code", err)
			} else if codeAnalysis != nil {
				response.CodeAnalysis = codeAnalysis
			}
		}
	}

	return response, nil
}

