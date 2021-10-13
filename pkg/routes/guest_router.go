package routes

import (
	"errors"
	"fmt"
	"github.com/LucasCarioca/go-template/pkg/config"
	"github.com/LucasCarioca/go-template/pkg/datasource"
	"github.com/LucasCarioca/go-template/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
)

//GuestRouter router for guest CRUD operations
type GuestRouter struct {
	db     *gorm.DB
	config *viper.Viper
}

//CreateGuestRequest structure of the create request for guests
type CreateGuestRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	RegistrationKey string `json:"registration_key" binding:"required"`
}

//NewGuestRouter creates a new instance of the guest router
func NewGuestRouter(app *gin.Engine) {
	r := GuestRouter{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}

	app.GET("/api/v1/guests", r.getAllGuests)
	app.GET("/api/v1/guests/:id", r.getGuest)
	app.POST("/api/v1/guests", r.createGuest)
	app.DELETE("/api/v1/guests/:id", r.deleteGuest)
}

func (r *GuestRouter) checkKey(ctx *gin.Context) error {
	apiKey := r.config.GetString("API_KEY")
	requestKey := ctx.Query("api_key")

	if apiKey != requestKey {
		return errors.New("INVALID_API_KEY")
	}
	return nil
}

func (r *GuestRouter) checkInvitation(ctx *gin.Context) (models.Invitation, error) {
	requestKey := ctx.Query("registration_key")
	i := models.Invitation{}
	if requestKey != "" {
		var c int64
		r.db.Where("registration_key = ?", requestKey).First(&i).Count(&c)
		if c < 1 {
			return i, errors.New("INVALID_REGISTRATION_KEY")
		}
		return i, nil
	}
	return i, errors.New("MISSING_REGISTRATION_KEY")
}

func (r *GuestRouter) readID(ctx *gin.Context) *int {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid id",
		})
		return nil
	}
	return &id
}

func (r *GuestRouter) getAllGuests(ctx *gin.Context) {
	guests := make([]models.Guest, 0)
	i, err := r.checkInvitation(ctx)
	if err != nil {
		err := r.checkKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
			return
		}
		r.db.Preload(clause.Associations).Find(&guests)
	} else {
		r.db.Table("guests").Where("invitation_id = ?", i.ID).Preload(clause.Associations).Find(&guests)
	}
	ctx.JSON(http.StatusOK, guests)
}

func (r *GuestRouter) createGuest(ctx *gin.Context) {
	var data CreateGuestRequest
	ctx.BindJSON(&data)
	fmt.Println(data)
	i, err := r.checkInvitation(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "INVITATION_NOT_FOUND"})
		return
	}

	var gc int64
	r.db.Table("guests").Where("invitation_id = ?", i.ID).Count(&gc)

	if gc >= int64(i.GuestCount) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invitation guest count limit reached", "error": "GUEST_COUNT_LIMIT"})
		return
	}

	g := &models.Guest{
		FirstName:  data.FirstName,
		LastName:   data.LastName,
		Email:      data.Email,
		Phone:      data.Phone,
		Invitation: i,
	}
	r.db.Create(g)
	ctx.JSON(http.StatusOK, g)
}

func (r *GuestRouter) getGuest(ctx *gin.Context) {
	id := r.readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}

	g, err := r.findGuestByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": err.Error()})
		return
	}

	i, err := r.checkInvitation(ctx)
	if err != nil {
		err = r.checkKey(ctx)
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

func (r *GuestRouter) findGuestByID(id int) (*models.Guest, error) {
	g := models.Guest{}
	var c int64
	r.db.Preload(clause.Associations).Find(&g, id).Count(&c)
	if c > 0 {
		return &g, nil
	}
	return nil, errors.New("GUEST_NOT_FOUND")
}

func (r *GuestRouter) deleteGuest(ctx *gin.Context) {
	id := r.readID(ctx)
	if id == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}

	g, err := r.findGuestByID(*id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		return
	}

	i, err := r.checkInvitation(ctx)
	if err != nil {
		err = r.checkKey(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request", "error": err.Error()})
			return
		}
		r.db.Delete(&g)
		ctx.JSON(http.StatusOK, g)
		return
	}

	if i.ID == g.InvitationID {
		r.db.Delete(&g)
		ctx.JSON(http.StatusOK, g)
		return
	}

	ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
}
