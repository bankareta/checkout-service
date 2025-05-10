package config

import (
	"base-golang/internal/delivery/http/route"
	"base-golang/internal/environtment"

	httpEtalase "base-golang/internal/m_etalase/delivery/http"
	mapperEtalase "base-golang/internal/m_etalase/mapper"
	repositoryEtalase "base-golang/internal/m_etalase/repository"
	usecaseEtalase "base-golang/internal/m_etalase/usecase"

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
