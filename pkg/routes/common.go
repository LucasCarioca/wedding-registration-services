package routes

import (
	"errors"
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func checkKey(ctx *gin.Context) error {
	config := config.GetConfig()
	apiKey := config.GetString("API_KEY")
	requestKey := ctx.Query("api_key")
	if apiKey != requestKey {
		return errors.New("INVALID_API_KEY")
	}
	return nil
}

func readID(ctx *gin.Context) *int {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid id",
		})
		return nil
	}
	return &id
}
