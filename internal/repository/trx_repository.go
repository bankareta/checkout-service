package repository

import (
	"log"

	"gorm.io/gorm"
)

// GormTransactionManager implements TransactionManager for GORM.
type GormTransactionManager struct {
	db *gorm.DB
}

// NewGormTransactionManager creates a new instance of GormTransactionManager.
func NewGormTransactionManager(db *gorm.DB) TransactionManager {
	return &GormTransactionManager{
		db: db,
	}
}

// Execute runs the provided function within a transaction.
func (tm *GormTransactionManager) Execute(txFunc func(tx *gorm.DB) error) (any, error) {
	tx := tm.db.Begin()

	if tx.Error != nil {
		log.Println("Failed to start transaction:", tx.Error)
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to panic:", r)
			panic(r)
		}
	}()

	// Execute the function in a transaction
	if err := txFunc(tx); err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			log.Println("Failed to rollback transaction:", rollbackErr)
		}
		return nil, err
	}

	return nil, tx.Commit().Error
}

// TransactionManager defines the behavior of a transaction manager.
type TransactionManager interface {
	// Execute runs a function within a transaction.
	Execute(txFunc func(tx *gorm.DB) error) (any, error)
}
