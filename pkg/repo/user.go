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

// UserFilter for filtering users in pagination
type UserFilter struct {
	dto.BasePaginationQuery
	Email        *string   `json:"email,omitempty" form:"email"`
	FullName     *string   `json:"full_name,omitempty" form:"full_name"`
	GuardName    *string   `json:"guard_name,omitempty" form:"guard_name"`
	StartDate    *string   `json:"start_date,omitempty" form:"start_date"` // Format: YYYY-MM-DD
	EndDate      *string   `json:"end_date,omitempty" form:"end_date"`     // Format: YYYY-MM-DD
	RoleIDs      []uint64  `json:"role_ids,omitempty" form:"role_ids"`
	PermissionIDs []uint64 `json:"permission_ids,omitempty" form:"permission_ids"`
	GetFields    []string  `json:"get_fields,omitempty" form:"get_fields"`
}

type UserRepo interface {
	creator[*model.User]
	reader[*model.User, *UserFilter]
	getByID[*model.User]
	updater[*model.User]
	GetList(ctx context.Context) ([]model.User, error)
	GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.User, error)
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

func (s *UserRepoImpl) Paginate(fields *UserFilter) (*dto.PaginationResponse[*model.User], error) {
	var items []*model.User
	var total int64
	query := s.db.Model(&model.User{}).Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

	// Filter by email
	if fields.Email != nil && *fields.Email != "" {
		query = query.Where("users.email LIKE ?", "%"+*fields.Email+"%")
	}

	// Filter by full_name
	if fields.FullName != nil && *fields.FullName != "" {
		query = query.Where("users.full_name LIKE ?", "%"+*fields.FullName+"%")
	}

	// Filter by guard_name
	if fields.GuardName != nil && *fields.GuardName != "" {
		query = query.Where("users.guard_name = ?", *fields.GuardName)
	}

	// Filter by date range
	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	if fields.StartDate != nil && *fields.StartDate != "" {
		startTime, err := time.ParseInLocation(layout, *fields.StartDate, loc)
		if err == nil {
			query = query.Where("users.created_at >= ?", startTime)
		}
	}

	if fields.EndDate != nil && *fields.EndDate != "" {
		endTime, err := time.ParseInLocation(layout, *fields.EndDate, loc)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("users.created_at <= ?", endTime)
		}
	}

	// Filter by roles
	if len(fields.RoleIDs) > 0 {
		query = query.Joins("JOIN model_has_roles ON users.id = model_has_roles.model_id").
			Where("model_has_roles.model_type = ?", "User").
			Where("model_has_roles.role_id IN ?", fields.RoleIDs).
			Group("users.id")
	}

	// Filter by permissions (direct or via roles)
	if len(fields.PermissionIDs) > 0 {
		// Get user IDs with direct permissions
		var directUserIDs []uint64
		s.db.Model(&model.ModelHasPermission{}).
			Joins("JOIN users ON model_has_permissions.model_id = users.id").
			Where("model_has_permissions.model_type = ?", "User").
			Where("model_has_permissions.permission_id IN ?", fields.PermissionIDs).
			Pluck("users.id", &directUserIDs)

		// Get user IDs with permissions via roles
		var roleUserIDs []uint64
		s.db.Model(&model.ModelHasRole{}).
			Joins("JOIN role_has_permissions ON model_has_roles.role_id = role_has_permissions.role_id").
			Joins("JOIN users ON model_has_roles.model_id = users.id").
			Where("model_has_roles.model_type = ?", "User").
			Where("role_has_permissions.permission_id IN ?", fields.PermissionIDs).
			Pluck("users.id", &roleUserIDs)

		// Combine both lists
		allUserIDs := make(map[uint64]bool)
		for _, id := range directUserIDs {
			allUserIDs[id] = true
		}
		for _, id := range roleUserIDs {
			allUserIDs[id] = true
		}

		// Convert map to slice
		userIDs := make([]uint64, 0, len(allUserIDs))
		for id := range allUserIDs {
			userIDs = append(userIDs, id)
		}

		if len(userIDs) > 0 {
			query = query.Where("users.id IN ?", userIDs)
		} else {
			// No users found with these permissions
			query = query.Where("1 = 0")
		}
	}

	// Preload relations
	query = query.Preload("Roles").Preload("Permissions")

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	// Select specific fields if provided
	if len(fields.GetFields) != 0 {
		query = query.Select(fields.GetFields)
	}

	// Pagination
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
	if totalPage == 0 {
		totalPage = 1
	}
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
		Page:      page,
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}

func (s *UserRepoImpl) UpdateSelectField(entity *model.User, fields map[string]interface{}) error {
	return s.db.Model(entity).Updates(fields).Error
}
