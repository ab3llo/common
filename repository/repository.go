package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, model *T) (*T, error)
	Get(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context, limit, offset int) ([]T, error)
	GetAllByField(ctx context.Context, fieldName, fieldValue string) ([]T, error)
	Update(ctx context.Context, id string, model *T) (*T, error)
	Delete(ctx context.Context, id string) error
}

type repository[T any] struct {
	DB *gorm.DB
}

func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &repository[T]{DB: db}
}

func (r *repository[T]) Create(ctx context.Context, model *T) (*T, error) {
	if err := r.DB.Create(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *repository[T]) Get(ctx context.Context, id string) (*T, error) {
	var model T
	if err := r.DB.First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &model, nil
}
func (r *repository[T]) GetAll(ctx context.Context, limit, offset int) ([]T, error) {
	var models []T
	if err := r.DB.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return models, nil
}

func (r *repository[T]) GetAllByField(ctx context.Context, fieldName, fieldValue string) ([]T, error) {
	var models []T
	var query = fmt.Sprintf("%s = ?", fieldName)
	if err := r.DB.Find(&models, query, fieldValue).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return models, nil
}

func (r *repository[T]) Update(ctx context.Context, id string, model *T) (*T, error) {
	var existing T
	if err := r.DB.First(&existing, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	if err := r.DB.Save(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *repository[T]) Delete(ctx context.Context, id string) error {
	var existing T
	if err := r.DB.Delete(&existing, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}
