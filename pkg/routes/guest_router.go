package routes

import (
	"net/http"

	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/LucasCarioca/wedding-registration-services/pkg/models"
	"github.com/LucasCarioca/wedding-registration-services/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//GuestRouter router for guest CRUD operations
type GuestRouter struct {
	db     *gorm.DB
	config *viper.Viper
	gs     *services.GuestService
	is     *services.InvitationService
}

//CreateGuestRequest structure of the create request for guests
type CreateGuestRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
}

//NewGuestRouter creates a new instance of the guest router
func NewGuestRouter(app *gin.Engine) {
	r := GuestRouter{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
		gs:     services.NewGuestService(),
		is:     services.NewInvitationService(),
	}

	app.GET("/api/v1/guests", r.getAllGuests)
	app.GET("/api/v1/guests/:id", r.getGuest)
	app.POST("/api/v1/guests", r.createGuest)
	app.DELETE("/api/v1/guests/:id", r.deleteGuest)
}

func (r *GuestRouter) checkInvitation(ctx *gin.Context) (*models.Invitation, error) {
	requestKey := ctx.Query("registration_key")
	return r.is.GetInvitationByRegistrationKey(requestKey)
}

func (r *GuestRouter) getAllGuests(ctx *gin.Context) {
	guests := make([]models.Guest, 0)
	i, err := r.checkInvitation(ctx)
	if err != nil {
		err := checkKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
			return
		}
		guests = r.gs.GetAllGuests()
	} else {
		guests, err = r.gs.GetAllGuestsByInvitationID(i.ID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "guests not found", "error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, guests)
}

func (r *GuestRouter) createGuest(ctx *gin.Context) {
	var data CreateGuestRequest
	ctx.BindJSON(&data)
	i, err := r.checkInvitation(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND"})
		return
	}

	if r.gs.GetGuestCountByInvitationID(i.ID) >= i.GuestCount {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invitation guest count limit reached", "error": "GUEST_COUNT_LIMIT"})
		return
	}
	g := r.gs.CreateGuest(data.FirstName, data.LastName, data.Email, data.Phone, *i)
	ctx.JSON(http.StatusOK, g)
}

func (r *GuestRouter) getGuest(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}

	g, err := r.gs.GetGuestByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": err.Error()})
		return
	}

	i, err := r.checkInvitation(ctx)
	if err != nil {
		err = checkKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, g)
		return
	}

	if i.ID == g.InvitationID {
		ctx.JSON(http.StatusOK, g)
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
}

func (r *GuestRouter) deleteGuest(ctx *gin.Context) {
	id := readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}

	i, err := r.checkInvitation(ctx)
	if err != nil {
		err = checkKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
			return
		}
		g, err := r.gs.DeleteGuestByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
			return
		}
		ctx.JSON(http.StatusOK, g)
		return
	}

	g, err := r.gs.GetGuestByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}
	if i.ID == g.InvitationID {
		r.gs.DeleteGuestByID(*id)
		ctx.JSON(http.StatusOK, g)
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
}
