package repository

import (
	"checkout-service/internal/entity"
	"checkout-service/internal/helpers"
	repoGlobal "checkout-service/internal/repository"
	tr_transaction "checkout-service/internal/tr_transaction"
	"context"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type transactionRepository struct {
	repoGlobal.Repository[entity.Transaction]
	TransactionDetailRepo repoGlobal.Repository[entity.TransactionDetail]
	Log                   *logrus.Logger
}

func NewTransactionRepository(db *gorm.DB, log *logrus.Logger) tr_transaction.TransactionRepository {
	return &transactionRepository{
		Log: log,
		Repository: repoGlobal.Repository[entity.Transaction]{
			DB: db,
		},
		TransactionDetailRepo: repoGlobal.Repository[entity.TransactionDetail]{
			DB: db,
		},
	}
}

func (r *transactionRepository) ExecuteTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
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

func (r *transactionRepository) ExecuteTransactionWithResult(ctx context.Context, fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	result, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return result, tx.Commit().Error
}

func (r *transactionRepository) AddTransaction(tx *gorm.DB, trans entity.Transaction) (entity.Transaction, error) {
	ctx := context.Background()
	if err := helpers.NewTrxManager(r.DB).WithTrx(&ctx, func(ctx context.Context) error {
		if err := tx.Table("transaction").Create(&trans).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return trans, err
	}
	return trans, nil
}

func (r *transactionRepository) AddDetailTransaction(tx *gorm.DB, transDetail entity.TransactionDetail) error {
	ctx := context.Background()
	if err := helpers.NewTrxManager(r.DB).WithTrx(&ctx, func(ctx context.Context) error {
		if err := tx.Table("transaction_detail").Create(&transDetail).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
