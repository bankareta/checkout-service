package repository

import (
	"checkout-service/internal/entity"
	master_products "checkout-service/internal/master_products"
	repoGlobal "checkout-service/internal/repository"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type masterProductsRepository struct {
	repoGlobal.Repository[entity.Products]
	Log *logrus.Logger
}

func NewMasterProductsRepository(db *gorm.DB, log *logrus.Logger) master_products.MasterProductsRepository {
	return &masterProductsRepository{
		Log: log,
		Repository: repoGlobal.Repository[entity.Products]{
			DB: db,
		},
	}
}

func (r *masterProductsRepository) ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *masterProductsRepository) GetProducts() ([]entity.Products, error) {
	var result []entity.Products
	query := r.DB.Table("products").Where("deleted_at IS NULL")
	err := query.Order(fmt.Sprintf("%s %s", "id", "DESC")).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *masterProductsRepository) ScanProduct(sku string) (*entity.Products, error) {
	var entity entity.Products
	result := r.DB.Table("products").Where("sku = ? AND deleted_at IS NULL", sku).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &entity, nil
}

func (r *masterProductsRepository) DetailProduct(id uint) (*entity.Products, error) {
	var entity entity.Products
	result := r.DB.Table("products").Where("id = ? AND deleted_at IS NULL", id).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &entity, nil
}
