package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
)

type ReqCoaAccountRepo interface {
	CountByFields(ctx context.Context, fields map[string]interface{}) (int64, error)
}

type reqCoaAccountRepo struct {
	db *gorm.DB
}

func NewReqCoaAccountRepo(db *gorm.DB) ReqCoaAccountRepo {
	return &reqCoaAccountRepo{
		db: db,
	}
}

func (s *reqCoaAccountRepo) CountByFields(ctx context.Context, fields map[string]interface{}) (int64, error) {
	var count int64

	err := s.db.WithContext(ctx).
		Model(&model.RequestCoaAccount{}).
		Where(fields).
		Count(&count).Error

	return count, err
}
