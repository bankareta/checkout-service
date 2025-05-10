package main

import (
	"checkout-service/internal/config"
	"checkout-service/internal/environtment"
	"fmt"
)

func init() {

	environtment.LoadConfig()
}

func main() {
	configEnv := environtment.Configs
	log := config.NewLogger(configEnv)
	db := config.NewDatabase(configEnv, log)
	validate := config.NewValidator(configEnv)
	app := config.NewFiber(configEnv)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   configEnv,
	})

	webPort := configEnv.WEB_PORT
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
