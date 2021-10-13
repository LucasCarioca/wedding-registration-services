package routes

import (
	"errors"
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/LucasCarioca/wedding-registration-services/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

//InvitationRouter  router for invitation CRUD operations
type InvitationRouter struct {
	db     *gorm.DB
	config *viper.Viper
}

//CreateInvitationRequest structure of the create request for invitations
type CreateInvitationRequest struct {
	Name       string `json:"name" binding:"required"`
	GuestCount int    `json:"guest_count" binding:"required"`
}

//NewInvitationRouter creates a new instance of the invitation router
func NewInvitationRouter(app *gin.Engine) {
	r := InvitationRouter{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}

	app.GET("/api/v1/invitations", r.getAllInvitations)
	app.GET("/api/v1/invitations/:id", r.getInvitation)
	app.POST("/api/v1/invitations", r.createInvitation)
	app.DELETE("/api/v1/invitations/:id", r.deleteInvitation)
}

func (r *InvitationRouter) checkKey(ctx *gin.Context) error {
	apiKey := r.config.GetString("API_KEY")
	requestKey := ctx.Query("api_key")

	if apiKey != requestKey {
		return errors.New("INVALID_API_KEY")
	}
	return nil
}

func (r *InvitationRouter) readID(ctx *gin.Context) *int {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid id",
		})
		return nil
	}
	return &id
}

func (r *InvitationRouter) getAllInvitations(ctx *gin.Context) {
	key := ctx.Query("registration_key")
	if key != "" {
		i := models.Invitation{}
		var c int64
		r.db.Where("registration_key = ?", key).First(&i).Count(&c)
		if c < 1 {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusOK, i)
		return
	}

	err := r.checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}

	invitations := make([]models.Invitation, 0)
	r.db.Find(&invitations)
	ctx.JSON(http.StatusOK, invitations)
}

func (r *InvitationRouter) createInvitation(ctx *gin.Context) {
	err := r.checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	var data CreateInvitationRequest
	ctx.BindJSON(&data)
	i := &models.Invitation{
		Name:       data.Name,
		GuestCount: data.GuestCount,
		Registered: false,
	}
	r.db.Create(i)
	ctx.JSON(http.StatusOK, i)
}

func (r *InvitationRouter) getInvitation(ctx *gin.Context) {
	err := r.checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	id := r.readID(ctx)
	if id != nil {
		i := models.Invitation{}
		var c int64
		r.db.Find(&i, id).Count(&c)
		if c > 0 {
			ctx.JSON(http.StatusOK, i)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{})
		}
	}
}

func (r *InvitationRouter) deleteInvitation(ctx *gin.Context) {
	err := r.checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	id := r.readID(ctx)
	if id != nil {
		i := models.Invitation{}
		var c int64
		r.db.Find(&i, id).Count(&c)
		if c > 0 {
			r.db.Delete(&i)
			ctx.JSON(http.StatusOK, i)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{})
		}
	}
}
