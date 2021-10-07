package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var config *viper.Viper

func Init(env string) {
	config = viper.New()

	config.SetConfigName(fmt.Sprintf("config.%s", env))
	config.SetConfigType("yaml")
	config.AddConfigPath(".")
	config.SetEnvPrefix("secret")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	for _, k := range config.AllKeys() {
		value := config.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			config.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value,"${"), "}")))
		}
	}
}

func GetConfig() *viper.Viper { 
	return config
}

func getEnvOrPanic(raw string) string {
	values := strings.Split(raw, "|")
	env := values[0]
	res := os.Getenv(env)
	if len(res) == 0 {
		if len(values) == 2 {
			return values[1]
		}
		panic("Mandatory env variable not found:" + env)
	}
	return res
}