package config

import (
	"checkout-service/internal/environtment"

	"github.com/sirupsen/logrus"
)

func NewLogger(viper *environtment.Config) *logrus.Logger {
	log := logrus.New()

	log.SetLevel(logrus.Level(viper.LOG_LEVEL))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
