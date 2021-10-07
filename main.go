package main

import (
	"flag"
	"fmt"
	"github.com/LucasCarioca/go-template/pkg/config"
	"github.com/LucasCarioca/go-template/pkg/datasource"
	"github.com/LucasCarioca/go-template/pkg/server"
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
	datasource.Init()
	server.Init(config.GetConfig())
}
