package ports

import (
	"context"
)

type (
	Pagination[T interface{}] struct {
		Data  []T
		Count int64
	}
	Repository[T interface{}] interface {
		GetAll(ctx context.Context,
			filter map[string]interface{}) []T
		GetAllSkipTake(ctx context.Context,
			filter map[string]interface{},
			skip int64,
			take int64) *Pagination[T]
		Count(ctx context.Context,
			filter map[string]interface{}) int64
		GetFirst(ctx context.Context,
			filter map[string]interface{}) *T
		Insert(ctx context.Context,
			entity *T)
		InsertAll(ctx context.Context,
			entities []T)
		Replace(ctx context.Context,
			filter map[string]interface{},
			entity *T)
		DeleteAll(ctx context.Context,
			filter map[string]interface{})
		Aggregate(
			ctx context.Context,
			pipeline []map[string]interface{}) []map[string]interface{}
	}
)
