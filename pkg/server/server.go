package server

import (
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"github.com/LucasCarioca/go-template/pkg/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"fmt"
)

func routesInit(app *gin.Engine) {
	routes.HealthRouter(app)
}

func Init(config *viper.Viper) {
	app := gin.Default()
	app.Use(cors.Default())
	app.Use(static.Serve("/", static.LocalFile(config.GetString("server.static"), false)))
	routesInit(app)
	host := config.GetString("server.host")
	port := config.GetString("server.port")
	app.Run(fmt.Sprintf("%s:%s", host, port))
}