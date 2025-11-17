package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"fmt"

	"core-ledger/pkg/utils/helper"
	wv "core-ledger/pkg/utils/wrapvalue"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CoAccountRepo interface {
	creator[*model.CoaAccount]
	reader[*model.CoaAccount, *dto.ListCoaAccountFilter]
	getByID[*model.CoaAccount]
	updater[*model.CoaAccount]

	Save(customer *model.CoaAccount) error
	Upsert(accounts []*model.CoaAccount, updateColumns []string) error
	GetParentID(ctx context.Context, id string) (*uint64, error)
	PaginateWithScopes(ctx context.Context, filter *dto.ListCoaAccountFilter) (*dto.PaginationResponse[*model.CoaAccount], error)
}

type coAccountRepo struct {
	db *gorm.DB
}

func NewCoAccountRepo(db *gorm.DB) CoAccountRepo {
	return &coAccountRepo{
		db: db,
	}
}
func (c *coAccountRepo) Save(customer *model.CoaAccount) error {
	return c.db.Create(&customer).Error
}

func (c *coAccountRepo) Create(customer ...*model.CoaAccount) error {
	return c.db.Create(customer).Error
}

func (c *coAccountRepo) GetByID(ctx context.Context, id int64) (*model.CoaAccount, error) {
	account := &model.CoaAccount{}

	err := c.db.Preload("Entries.Journal").First(&account, id).Error
	if err != nil {
		return nil, err
	}
	m := make(map[uint64]model.Journal)
	for _, e := range account.Entries {
		m[e.JournalID] = *e.Journal
	}

	account.Journals = make([]model.Journal, 0, len(m))
	for _, j := range m {
		account.Journals = append(account.Journals, j)
	}
	return account, err
}

func (c *coAccountRepo) Update(customer *model.CoaAccount) error {
	return c.db.Save(&customer).Error
}

func (c *coAccountRepo) UpdateSelectField(entity *model.CoaAccount, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *coAccountRepo) Upsert(accounts []*model.CoaAccount, updateColumns []string) error {
	if len(accounts) == 0 {
		return nil
	}

	if len(updateColumns) == 0 {
		// Mặc định update tất cả trường có thể thay đổi
		updateColumns = []string{"name", "type", "parent_id", "status", "provider", "network", "tags", "metadata", "updated_at"}
	}
	for _, account := range accounts {
		// account_no tạo khi chưa có
		if account.AccountNo == "" {
			account.AccountNo = helper.GenerateSecureNumber()
		}
	}
	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "currency"}}, // cột unique
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&accounts).Error
}

func (c *coAccountRepo) GetParentID(ctx context.Context, parentCode string) (*uint64, error) {
	if parentCode == "" {
		return nil, nil
	}

	var parent model.CoaAccount
	if err := c.db.WithContext(ctx).Select("id").Where("code = ?", parentCode).Where("parent_id IS NULL").First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // không tìm thấy → trả nil
		}
		return nil, err // lỗi DB khác
	}

	return &parent.ID, nil
}

func (s *coAccountRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.CoaAccount, error) {
	var coaAccount *model.CoaAccount
	query := s.db.WithContext(ctx).Model(&model.CoaAccount{})
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return coaAccount, query.Where(fields).First(&coaAccount).Error
}

func (t *coAccountRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.CoaAccount, error) {
	var coa []*model.CoaAccount
	query := t.db.WithContext(ctx).Model(&model.CoaAccount{}).Where(fields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return coa, query.Find(&coa).Error
}

func (s *coAccountRepo) Paginate(fields *dto.ListCoaAccountFilter) (*dto.PaginationResponse[*model.CoaAccount], error) {

	var total int64
	fmt.Println("Status", fields.Status)
	query := s.db.Model(&model.CoaAccount{}).Order("id DESC")
	if fields.Search != nil && *fields.Search != "" {
		likeQuery := "%" + *fields.Search + "%"
		query = query.Where(
			s.db.
				Where("name LIKE ?", likeQuery).
				Or("code LIKE ?", likeQuery).
				Or("account_no LIKE ?", likeQuery),
		)
	}

	if len(fields.Types) > 0 {
		query = query.Where("type IN (?)", fields.Types)
	}
	if len(fields.Status) > 0 {
		query = query.Where("status IN (?)", fields.Status)
	}
	if len(fields.Networks) > 0 {
		query = query.Where("network IN (?)", fields.Networks)
	}
	if len(fields.Providers) > 0 {
		query = query.Where("provider IN (?)", fields.Providers)
	}
	err := query.Count(&total).Error
	items := make([]*model.CoaAccount, 0, total)
	if err != nil {
		return nil, err
	}

	limit := 25
	offset := 0
	var page int64 = 1
	if fields.Limit != nil {
		limit = int(*fields.Limit)
	}
	if fields.Page != nil {
		offset = int(*fields.Page-1) * limit
		page = *fields.Page
	}

	err = query.Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, err
	}
	totalPage := total/int64(limit) + 1
	var nextPage *int64
	var prevPage *int64
	if page < totalPage {
		nextPage = wv.ToPointer(page + 1)
	}
	if page > 1 {
		prevPage = wv.ToPointer(page - 1)
	}
	return &dto.PaginationResponse[*model.CoaAccount]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}
func (r *coAccountRepo) PaginateWithScopes(ctx context.Context, fields *dto.ListCoaAccountFilter) (*dto.PaginationResponse[*model.CoaAccount], error) {

	params := BuildParamsFromFilter(fields)

	var items []*model.CoaAccount

	// q = q.
	// 	Preload("Entries").Preload("Parent").
	// 	Preload("Children")
	limit := int64(25)
	page := int64(1)
	if fields.Limit != nil {
		limit = *fields.Limit
	}
	if fields.Page != nil {
		page = *fields.Page
	}

	pagination, err := CustomPaginate(r.db.Model(&model.CoaAccount{}), params, page, limit, &items)
	if err != nil {
		return nil, err
	}

	return pagination, nil
	// if fields == nil {
	// 	fields = &dto.ListCoaAccountFilter{}
	// }

	// // Thiết lập giá trị mặc định
	// limit := int64(25)
	// page := int64(1)
	// if fields.Limit != nil && *fields.Limit > 0 {
	// 	limit = *fields.Limit
	// }
	// if fields.Page != nil && *fields.Page > 0 {
	// 	page = *fields.Page
	// }

	// // Tạo model instance để gọi scope methods
	// coaAccount := &model.CoaAccount{}

	// // Build query với context và áp dụng các scope methods
	// query := r.db.WithContext(ctx).Model(&model.CoaAccount{})

	// // Áp dụng scope search
	// if fields.Search != nil && *fields.Search != "" {
	// 	query = query.Scopes(coaAccount.ScopeSearch(*fields.Search))
	// }

	// // Áp dụng scope status
	// if len(fields.Status) > 0 {
	// 	query = query.Scopes(coaAccount.ScopeStatus(fields.Status))
	// }

	// // Áp dụng scope types (convert []*string sang []string)
	// if len(fields.Types) > 0 {
	// 	types := make([]string, 0, len(fields.Types))
	// 	for _, t := range fields.Types {
	// 		if t != nil {
	// 			types = append(types, *t)
	// 		}
	// 	}
	// 	if len(types) > 0 {
	// 		query = query.Scopes(coaAccount.ScopeTypes(types))
	// 	}
	// }

	// // Áp dụng scope providers (convert []*string sang []string)
	// if len(fields.Providers) > 0 {
	// 	providers := make([]string, 0, len(fields.Providers))
	// 	for _, p := range fields.Providers {
	// 		if p != nil {
	// 			providers = append(providers, *p)
	// 		}
	// 	}
	// 	if len(providers) > 0 {
	// 		query = query.Scopes(coaAccount.ScopeProviders(providers))
	// 	}
	// }

	// // Áp dụng scope sort
	// if fields.Sort != nil && *fields.Sort != "" {
	// 	query = query.Scopes(coaAccount.ScopeSort(*fields.Sort))
	// }

	// // Filter networks nếu có (chưa có scope method, dùng trực tiếp)
	// if len(fields.Networks) > 0 {
	// 	networks := make([]string, 0, len(fields.Networks))
	// 	for _, n := range fields.Networks {
	// 		if n != nil {
	// 			networks = append(networks, *n)
	// 		}
	// 	}
	// 	if len(networks) > 0 {
	// 		query = query.Where("network IN ?", networks)
	// 	}
	// }

	// // Filter currency nếu có (chưa có scope method, dùng trực tiếp)
	// if len(fields.Currency) > 0 {
	// 	currencies := make([]string, 0, len(fields.Currency))
	// 	for _, c := range fields.Currency {
	// 		if c != nil {
	// 			currencies = append(currencies, *c)
	// 		}
	// 	}
	// 	if len(currencies) > 0 {
	// 		query = query.Where("currency IN ?", currencies)
	// 	}
	// }

	// // Đếm tổng số bản ghi (không có sort và limit/offset)
	// var total int64
	// countQuery := query.Session(&gorm.Session{})
	// if err := countQuery.Count(&total).Error; err != nil {
	// 	return nil, fmt.Errorf("failed to count: %w", err)
	// }

	// // Query dữ liệu với pagination
	// var items []*model.CoaAccount
	// offset := int((page - 1) * limit)
	// if err := query.Limit(int(limit)).Offset(offset).Find(&items).Error; err != nil {
	// 	return nil, fmt.Errorf("failed to find: %w", err)
	// }

	// // Tính toán pagination
	// totalPage := (total + limit - 1) / limit
	// var nextPage, prevPage *int64
	// if page < totalPage {
	// 	n := page + 1
	// 	nextPage = &n
	// }
	// if page > 1 {
	// 	p := page - 1
	// 	prevPage = &p
	// }

	// return &dto.PaginationResponse[*model.CoaAccount]{
	// 	Items:     items,
	// 	Total:     total,
	// 	Limit:     limit,
	// 	Page:      page,
	// 	TotalPage: totalPage,
	// 	NextPage:  nextPage,
	// 	PrevPage:  prevPage,
	// }, nil
}
