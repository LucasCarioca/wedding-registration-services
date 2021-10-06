package config

import (
	"github.com/spf13/viper"
	"fmt"
)

var config *viper.Viper

func Init(env string) {
	config = viper.New()
	config.SetConfigName(fmt.Sprintf("config.%s", env)) 
	config.SetConfigType("yaml")
	config.AddConfigPath("config/")
	config.AddConfigPath("../config/")
	config.AddConfigPath("$HOME/config/")
	config.AddConfigPath("/config/")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func GetConfig() *viper.Viper { 
	return config
}