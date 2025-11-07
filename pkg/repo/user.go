package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	wv "core-ledger/pkg/utils/wrapvalue"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserRepo interface {
	creator[*model.User]
	reader[*model.User, *TransactionFilter]
	getByID[*model.User]
	// updater[*model.User]
	GetList(ctx context.Context) ([]model.User, error)
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoIml(db *gorm.DB) UserRepo {
	return &UserRepoImpl{db: db}
}
func (s *UserRepoImpl) Create(transactions ...*model.User) error {
	return s.db.Create(transactions).Error
}

func (s *UserRepoImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	customer := &model.User{}
	return customer, s.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (s *UserRepoImpl) Saves(trans []*model.User) error {
	return s.db.Model(&model.User{}).Save(trans).Error
}

func (s *UserRepoImpl) GetByIds(ids []int64) ([]*model.User, error) {
	var sessions []*model.User
	return sessions, s.db.Model(&model.User{}).Where("id in (?)", ids).Find(&sessions).Error
}
func (r *UserRepoImpl) GetList(ctx context.Context) ([]model.User, error) {
	var users []model.User
	result := r.db.WithContext(ctx).Find(&users)
	return users, result.Error
}
func (s *UserRepoImpl) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.User, error) {
	var customer *model.User
	query := s.db.WithContext(ctx).Model(&model.User{})
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return customer, query.Where(fields).First(&customer).Error
}
func (t *UserRepoImpl) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.User, error) {
	var trans []*model.User
	query := t.db.WithContext(ctx).Model(&model.User{}).Where(fields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return trans, query.Find(&trans).Error
}

func (s *UserRepoImpl) Paginate(fields *TransactionFilter) (*dto.PaginationResponse[*model.User], error) {
	var items []*model.User
	var total int64
	query := s.db.Model(&model.User{}).Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	if fields.StartDate != nil {
		startTime, err := time.ParseInLocation(layout, *fields.StartDate, loc)
		if err == nil {
			query = query.Where("users.created_at >= ?", startTime)
		}
	}

	if fields.EndDate != nil {
		endTime, err := time.ParseInLocation(layout, *fields.EndDate, loc)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("users.created_at <= ?", endTime)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	if len(fields.GetFields) != 0 {
		query = query.Select(fields.GetFields)
	}

	limit := 1000000
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

	return &dto.PaginationResponse[*model.User]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}
