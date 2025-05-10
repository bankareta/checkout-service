package usecase

import (
	"checkout-service/internal/helpers"
	metalase "checkout-service/internal/m_etalase"
	"checkout-service/internal/model"
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type metalaseUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	EtalaseRepository metalase.MEtalaseRepository
	mapper            metalase.MEtalaseMapper
}

func NewMEtalaseUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	etalaseRepository metalase.MEtalaseRepository,
	mapper metalase.MEtalaseMapper) metalase.MEtalaseUseCase {
	return &metalaseUseCase{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		EtalaseRepository: etalaseRepository,
		mapper:            mapper,
	}
}

func (c metalaseUseCase) AddEtalase(ctx context.Context, request *model.AddEtalaseRequest) error {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err.Error())
		return helpers.ValidateError
	}

	if c.mapper.ContainSpecialChars(request.EtalaseName) {
		c.Log.Warnf("Failed save etalase : %+v", "Special Char")
		return helpers.NewError(400, nil, "01", "Etalase Name contains disallowed special characters")
	}
	dataEtalase := c.mapper.MapAddEtalaseRequestToEntity(*request)
	_, err = c.EtalaseRepository.PostEtalase(dataEtalase)
	if err != nil {
		c.Log.Warnf("Failed save etalase : %+v", err.Error())
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1364 {
			return helpers.NewError(400, nil, "01", "Gagal menambahkan etalase, coba lagi.")
		}
		return err
	}
	return nil
}
