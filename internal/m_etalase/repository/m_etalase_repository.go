package repository

import (
	"context"

	"base-golang/internal/entity"
	"base-golang/internal/helpers"
	metalase "base-golang/internal/m_etalase"
	repoGlobal "base-golang/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type metalaseRepository struct {
	repoGlobal.Repository[entity.Etalase]
	etalase_product repoGlobal.Repository[entity.EtalaseProduct]
	Log             *logrus.Logger
}

func NewMEtalaseRepository(db *gorm.DB, log *logrus.Logger) metalase.MEtalaseRepository {
	return &metalaseRepository{
		Log: log,
		Repository: repoGlobal.Repository[entity.Etalase]{
			DB: db,
		},
		etalase_product: repoGlobal.Repository[entity.EtalaseProduct]{
			DB: db,
		},
	}
}

func (r *metalaseRepository) PostEtalase(dataAddEtalase entity.Etalase) (entity.Etalase, error) {
	ctx := context.Background()
	if err := helpers.NewTrxManager(r.DB).WithTrx(&ctx, func(ctx context.Context) error {
		if err := r.Create(r.DB, &dataAddEtalase); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return dataAddEtalase, err
	}

	return dataAddEtalase, nil
}
