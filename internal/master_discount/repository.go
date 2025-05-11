package master_discount

import (
	"checkout-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type MasterDiscountRepository interface {
	ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetDiscountProduct(productId int) ([]entity.ProductsDiscount, error)
	GetDiscount(id int) (*entity.Discount, error)
	GetFreeProductIDs() (map[int]bool, error)
}
