package environtment

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return config
}

var Configs *Config

type Config struct {
	APP_NAME               string `mapstructure:"APP_NAME"`
	APP_ENV                string `mapstructure:"APP_ENV"`
	WEB_PREFORK            bool   `mapstructure:"WEB_PREFORK"`
	WEB_PORT               int    `mapstructure:"WEB_PORT"`
	LOG_LEVEL              int    `mapstructure:"LOG_LEVEL"`
	DATABASE_USERNAME      string `mapstructure:"DATABASE_USERNAME"`
	DATABASE_PASSWORD      string `mapstructure:"DATABASE_PASSWORD"`
	DATABASE_HOST          string `mapstructure:"DATABASE_HOST"`
	DATABASE_PORT          int    `mapstructure:"DATABASE_PORT"`
	DATABASE_NAME          string `mapstructure:"DATABASE_NAME"`
	DATABASE_POOL_IDLE     int    `mapstructure:"DATABASE_POOL_IDLE"`
	DATABASE_POOL_MAX      int    `mapstructure:"DATABASE_POOL_MAX"`
	DATABASE_POOL_LIFETIME int    `mapstructure:"DATABASE_POOL_LIFETIME"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./../")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../..")

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("EROR")
		return
	}
	err = viper.Unmarshal(&config)
	Configs = &config
	return
}
