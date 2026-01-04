package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, page, perPage int) ([]T, int64, error) {
	var entities []T
	var total int64

	r.DB.WithContext(ctx).Model(new(T)).Count(&total)

	offset := (page - 1) * perPage
	err := r.DB.WithContext(ctx).Offset(offset).Limit(perPage).Find(&entities).Error

	return entities, total, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}