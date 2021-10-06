
package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func get(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func HealthRouter(router *gin.Engine) {
	router.GET("/api/v1/health", get)
}