package main

import (
	"flag"
	"fmt"
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"ggithub.com/LucasCarioca/wedding-registration-services/pkg/server"
	"os"
)

func getEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		envFlag := flag.String("e", "dev", "")
		flag.Usage = func() {
			fmt.Println("Usage: server -e {mode}")
			os.Exit(1)
		}
		flag.Parse()
		env = *envFlag
	}
	return env
}

func main() {
	config.Init(getEnv())
	datasource.Init(config.GetConfig())
	server.Init(config.GetConfig())
}
