package routes

import (
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/LucasCarioca/wedding-registration-services/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
)

//InvitationRouter  router for invitation CRUD operations
type InvitationRouter struct {
	s      *services.InvitationService
	db     *gorm.DB
	config *viper.Viper
}

//CreateInvitationRequest structure of the create request for invitations
type CreateInvitationRequest struct {
	Name       string `json:"name" binding:"required"`
	Message    string `json:"message" binding:"required"`
	GuestCount int    `json:"guest_count" binding:"required"`
}

//NewInvitationRouter creates a new instance of the invitation router
func NewInvitationRouter(app *gin.Engine) {
	r := InvitationRouter{
		s:      services.NewInvitationService(),
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}

	app.GET("/api/v1/invitations", r.getAllInvitations)
	app.GET("/api/v1/invitations/:id", r.getInvitation)
	app.POST("/api/v1/invitations", r.createInvitation)
	app.DELETE("/api/v1/invitations/:id", r.deleteInvitation)
}

func (r *InvitationRouter) getAllInvitations(ctx *gin.Context) {
	key := ctx.Query("registration_key")
	if key != "" {
		i, err := r.s.GetInvitationByRegistrationKey(key)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
		return
	}

	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, r.s.GetAllInvitations())
}

func (r *InvitationRouter) createInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	var data CreateInvitationRequest
	ctx.BindJSON(&data)
	i := r.s.CreateInvitation(data.Name, data.Message, data.GuestCount)
	ctx.JSON(http.StatusOK, i)
}

func (r *InvitationRouter) getInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	id := readID(ctx)
	if id != nil {
		i, err := r.s.GetInvitationByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}

func (r *InvitationRouter) deleteInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
		return
	}
	id := readID(ctx)
	if id != nil {
		i, err := r.s.DeleteInvitationByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}
