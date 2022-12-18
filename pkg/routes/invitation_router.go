package routes

import (
	"fmt"
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
	Phone      string `json:"phone" binding:"required"`
	Email      string `json:"email" binding:"required"`
	GuestCount int    `json:"guest_count" binding:"required"`
}

type DeclineInvitationRequest struct {
	Instruction string `json:"instruction" binding:"required"`
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
	app.PUT("/api/v1/invitations/:id", r.declineInvitation)
	app.DELETE("/api/v1/invitations/:id", r.deleteInvitation)
}

func (r *InvitationRouter) getAllInvitations(ctx *gin.Context) {
	value := ctx.Query("value")
	if value != "" {
		i, err := r.s.Search(value)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
		return
	}

	key := ctx.Query("registration_key")
	if key != "" {
		i, err := r.s.GetByRegistrationKey(key)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
		return
	}

	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": "UNAUTHORIZED_REQUEST", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, r.s.GetAll())
}

func (r *InvitationRouter) createInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": "UNAUTHORIZED_REQUEST", "details": err.Error()})
		return
	}
	var data CreateInvitationRequest
	if err = ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing or incorrect fields received", "error": "INVITATION_CREATE_PAYLOAD_INVALID", "details": err.Error()})
		return
	}
	fmt.Println(data)
	i, err := r.s.Create(data.Name, data.Message, data.Email, data.Phone, data.GuestCount)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "failed to create invitation", "error": "INVITATION_CREATE_FAILED", "details": err.Error()})
		return
	}
	fmt.Println(i)
	//ctx.JSON(http.StatusOK, i)
}

func (r *InvitationRouter) getInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": "UNAUTHORIZED_REQUEST", "details": err.Error()})
		return
	}

	id := readID(ctx)
	if id != nil {
		i, err := r.s.GetByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}

func (r *InvitationRouter) declineInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": "UNAUTHORIZED_REQUEST", "details": err.Error()})
		return
	}
	var data DeclineInvitationRequest
	if err = ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing or incorrect fields received", "error": "INVITATION_DECLINE_PAYLOAD_INVALID", "details": err.Error()})
		return
	}
	fmt.Println(data)
	id := readID(ctx)
	if id != nil {
		if data.Instruction == "declined" {
			i, err := r.s.DeclineById(*id)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, i)
			return
		}
		i, err := r.s.GetByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}

func (r *InvitationRouter) deleteInvitation(ctx *gin.Context) {
	err := checkKey(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": "UNAUTHORIZED_REQUEST", "details": err.Error()})
		return
	}
	id := readID(ctx)
	if id != nil {
		i, err := r.s.DeleteByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}
