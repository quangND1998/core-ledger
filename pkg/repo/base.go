package repo

import (
	"context"
	"core-ledger/model/dto"
	"fmt"
)

type PreloadInput struct {
	Model     string
	Condition *string
	Value     interface{}
}

type creator[T any] interface {
	Create(entity ...T) error
}

type reader[T, P any] interface {
	getOneByFields[T]
	getManyByFields[T]
	paginator[T, P]
}

type updater[T any] interface {
	UpdateSelectField(entity T, fields map[string]interface{}) error
}

type deleter[T any] interface {
	Delete(entity T) error
}

type getByID[T any] interface {
	GetByID(ctx context.Context, id int64) (T, error)
}

type getByUuid[T any] interface {
	GetByUuid(ctx context.Context, id string) (T, error)
}

type getOneByFields[T any] interface {
	GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (T, error)
}

type getManyByFields[T any] interface {
	GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]T, error)
}

type paginator[T, P any] interface {
	Paginate(fields P) (*dto.PaginationResponse[T], error)
}

func ExecutePaginate[T any]() (*dto.PaginationResponse[T], error) {
	fmt.Println("implement me")
	return nil, nil
}
