package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type BaseRepository[T any] interface {
	FindByID(ctx context.Context, id any) (*T, error)
	FindAll(ctx context.Context, page, limit int, conditions map[string]any, order string) ([]T, int64, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id any) error
	WithTx(tx *gorm.DB) BaseRepository[T]
	DB() *gorm.DB
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

func (r *baseRepository[T]) DB() *gorm.DB {
	return r.db
}

func (r *baseRepository[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: tx}
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindAll(ctx context.Context, page, limit int, conditions map[string]any, order string) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))

	// Apply filtering conditions
	for key, val := range conditions {
		if val != nil && val != "" {
			query = query.Where(key, val)
		}
	}

	// Count total matching records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if limit > 0 {
		offset := 0
		if page > 1 {
			offset = (page - 1) * limit
		}
		query = query.Offset(offset).Limit(limit)
	}

	// Order by
	if order != "" {
		query = query.Order(order)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepository[T]) Delete(ctx context.Context, id any) error {
	var entity T
	result := r.db.WithContext(ctx).Delete(&entity, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
