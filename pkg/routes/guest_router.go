package routes

import (
	"github.com/LucasCarioca/go-template/pkg/datasource"
	"github.com/LucasCarioca/go-template/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type GuestRouter struct {
	db *gorm.DB
}

type CreateGuestRequest struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	RegistrationKey string `json:"registration_key" binding:"required"`
}

func NewGuestRouter(app *gin.Engine) {
	r := GuestRouter{
		db: datasource.GetDataSource(),
	}

	app.GET("/api/v1/guests", r.getAllGuests)
	app.GET("/api/v1/guests/:id", r.getGuest)
	app.POST("/api/v1/guests", r.createGuest)
	app.DELETE("/api/v1/guests/:id", r.deleteGuest)
}

func (r *GuestRouter) readId(ctx *gin.Context) *int {
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
	r.db.Find(&guests)
	ctx.JSON(http.StatusOK, guests)
}

func (r *GuestRouter) createGuest(ctx *gin.Context) {
	var data CreateGuestRequest
	ctx.BindJSON(&data)

	i := models.Invitation{}
	var c int64
	r.db.Where("registration_key = ?", data.RegistrationKey).First(&i).Count(&c)

	if c < 1 {
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
	id := r.readId(ctx)
	if id != nil {
		i := models.Guest{}
		var c int64
		r.db.Find(&i, id).Count(&c)
		if c > 0 {
			ctx.JSON(http.StatusOK, i)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		}
	}
}

func (r *GuestRouter) deleteGuest(ctx *gin.Context) {
	id := r.readId(ctx)
	if id != nil {
		i := models.Guest{}
		var c int64
		r.db.Find(&i, id).Count(&c)
		if c > 0 {
			r.db.Delete(&i)
			ctx.JSON(http.StatusOK, i)
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "guest not found", "error": "GUEST_NOT_FOUND"})
		}
	}
}
