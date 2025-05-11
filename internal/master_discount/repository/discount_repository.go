package repository

import (
	"checkout-service/internal/entity"
	master_discount "checkout-service/internal/master_discount"
	repoGlobal "checkout-service/internal/repository"
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type masterDiscountRepository struct {
	repoGlobal.Repository[entity.Discount]
	products_discount repoGlobal.Repository[entity.ProductsDiscount]
	Log               *logrus.Logger
}

func NewMasterDiscountRepository(db *gorm.DB, log *logrus.Logger) master_discount.MasterDiscountRepository {
	return &masterDiscountRepository{
		Log: log,
		Repository: repoGlobal.Repository[entity.Discount]{
			DB: db,
		},
		products_discount: repoGlobal.Repository[entity.ProductsDiscount]{
			DB: db,
		},
	}
}

func (r *masterDiscountRepository) ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
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

func (r *masterDiscountRepository) GetDiscountProduct(productId int) ([]entity.ProductsDiscount, error) {
	var result []entity.ProductsDiscount
	query := r.DB.Table("products_discounts").Where("product_id = ?", productId)
	err := query.Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *masterDiscountRepository) GetDiscount(id int) (*entity.Discount, error) {
	var entity entity.Discount
	result := r.DB.Table("discounts").Where("id = ? AND deleted_at IS NULL", id).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &entity, nil
}

func (r *masterDiscountRepository) GetFreeProductIDs() (map[int]bool, error) {
	var discounts []entity.Discount
	err := r.DB.Table("discounts").Where("free_id_product IS NOT NULL AND deleted_at IS NULL").Find(&discounts).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int]bool)
	for _, d := range discounts {
		result[d.FreeIDProduct] = true
	}
	return result, nil
}
