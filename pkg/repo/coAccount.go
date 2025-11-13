package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"

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
	customer := &model.CoaAccount{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
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
	query := s.db.Model(&model.CoaAccount{}).Order("id asc")
	if fields.Name != nil && *fields.Name != "" {
		likeQuery := "%" + *fields.Name + "%"
		query = query.Where("name LIKE ?", likeQuery)
	}
	if fields.Code != nil && *fields.Code != "" {
		likeQuery := "%" + *fields.Code + "%"
		query = query.Where("code LIKE ?", likeQuery)
	}
	if fields.Status != nil && *fields.Status != "" {
		query = query.Where("status = ?", *fields.Status)
	}

	if fields.Type != nil && *fields.Type != "" {
		query = query.Where("type = ?", *fields.Type)
	}

	if fields.AccountNo != nil && *fields.AccountNo != "" {
		likeQuery := "%" + *fields.AccountNo + "%"
		query = query.Where("account_no LIKE ?", likeQuery)
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
