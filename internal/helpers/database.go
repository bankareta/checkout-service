package helpers

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type trxManager struct {
	db *gorm.DB
}

type trxFn func(ctx context.Context) error

func NewTrxManager(db *gorm.DB) *trxManager {
	return &trxManager{db}
}

func (g *trxManager) WithTrx(ctx *context.Context, fn trxFn) (err error) {
	tx := g.db.Begin()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logrus.Error(p)
			err = errors.New("panic happened because: " + fmt.Sprintf("%v", p))
		} else if err != nil {

			tx.Rollback()
		} else {

			err = tx.Commit().Error
		}
	}()

	err = fn(*ctx)
	return err
}
