package usecase

import (
	"checkout-service/internal/helpers"
	master_products "checkout-service/internal/master_products"
	"checkout-service/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type masterProductsUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	MasterProductsRepository master_products.MasterProductsRepository
	mapper                   master_products.MasterProductsMapper
}

func NewMasterProductsUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	masterProductsRepository master_products.MasterProductsRepository,
	mapper master_products.MasterProductsMapper) master_products.MasterProductsUseCase {
	return &masterProductsUseCase{
		DB:                       db,
		Log:                      logger,
		Validate:                 validate,
		MasterProductsRepository: masterProductsRepository,
		mapper:                   mapper,
	}
}

func (c masterProductsUseCase) InquiryProducts() ([]model.InquiryProductResp, error) {
	var res []model.InquiryProductResp

	resProducts, err := c.MasterProductsRepository.GetProducts()
	if err != nil {
		c.Log.Warnf("Failed Inquiry Products  : %+v", err)
		return res, helpers.NewError(500, nil, "500", err.Error())
	}

	respMap := c.mapper.MapInquiryProductsResp(res, resProducts)
	return respMap, nil
}
