package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func get(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// HealthRouter router for the health check endpoint
func HealthRouter(router *gin.Engine) {
	router.GET("/api/v1/health", get)
}
