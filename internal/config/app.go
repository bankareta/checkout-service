package config

import (
	"checkout-service/internal/delivery/http/route"
	"checkout-service/internal/environtment"

	httpMasterProducts "checkout-service/internal/master_products/delivery/http"
	mapperMasterProducts "checkout-service/internal/master_products/mapper"
	repositoryMasterProducts "checkout-service/internal/master_products/repository"
	usecaseMasterProducts "checkout-service/internal/master_products/usecase"

	repositoryMasterDiscount "checkout-service/internal/master_discount/repository"

	httpTransaction "checkout-service/internal/tr_transaction/delivery/http"
	mapperTransaction "checkout-service/internal/tr_transaction/mapper"
	repositoryTransaction "checkout-service/internal/tr_transaction/repository"
	usecaseTransaction "checkout-service/internal/tr_transaction/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *environtment.Config
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	masterProductsRepository := repositoryMasterProducts.NewMasterProductsRepository(config.DB, config.Log)
	masterDiscountRepository := repositoryMasterDiscount.NewMasterDiscountRepository(config.DB, config.Log)
	transactionRepository := repositoryTransaction.NewTransactionRepository(config.DB, config.Log)

	// setup mapper
	masterProductsMapper := mapperMasterProducts.NewMasterProductsMapper(config.Log)
	transactionMapper := mapperTransaction.NewTransactionMapper(config.Log)

	// setup use cases
	masterProductsUseCase := usecaseMasterProducts.NewMasterProductsUseCase(config.DB, config.Log, config.Validate, masterProductsRepository, masterProductsMapper)
	transactionUseCase := usecaseTransaction.NewTransactionUseCase(config.DB, config.Log, config.Validate, transactionRepository, masterProductsRepository, masterDiscountRepository, transactionMapper)

	// setup controller
	masterProductsController := httpMasterProducts.NewMasterProductsController(masterProductsUseCase, config.Log)
	transactionController := httpTransaction.NewTransactionController(transactionUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:                      config.App,
		MasterProductsController: masterProductsController,
		TransactionController:    transactionController,
	}
	routeConfig.Setup()
}
