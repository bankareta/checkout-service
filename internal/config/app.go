package config

import (
	"checkout-service/internal/delivery/http/route"
	"checkout-service/internal/environtment"

	httpEtalase "checkout-service/internal/m_etalase/delivery/http"
	mapperEtalase "checkout-service/internal/m_etalase/mapper"
	repositoryEtalase "checkout-service/internal/m_etalase/repository"
	usecaseEtalase "checkout-service/internal/m_etalase/usecase"

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
	metalaseRepository := repositoryEtalase.NewMEtalaseRepository(config.DB, config.Log)

	// setup mapper
	etalaseMapper := mapperEtalase.NewMEtalaseMapper(config.Log)

	// setup use cases
	etalaseUseCase := usecaseEtalase.NewMEtalaseUseCase(config.DB, config.Log, config.Validate, metalaseRepository, etalaseMapper)

	// setup controller
	etalaseController := httpEtalase.NewMEtalaseController(etalaseUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:               config.App,
		EtalaseController: etalaseController,
	}
	routeConfig.Setup()
}
