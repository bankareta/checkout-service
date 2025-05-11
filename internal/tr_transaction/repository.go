package tr_transaction

import (
	"checkout-service/internal/entity"
	"context"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	ExecuteTransactionWithResult(ctx context.Context, fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error)
	AddTransaction(tx *gorm.DB, trans entity.Transaction) (entity.Transaction, error)
	AddDetailTransaction(tx *gorm.DB, transDetail entity.TransactionDetail) error
}
