package master_products

import (
	"checkout-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type MasterProductsRepository interface {
	ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetProducts() ([]entity.Products, error)
	ScanProduct(sku string) (*entity.Products, error)
	DetailProduct(id uint) (*entity.Products, error)
}
