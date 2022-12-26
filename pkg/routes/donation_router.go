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

// DonationRouter  router for donation CRUD operations
type DonationRouter struct {
	s      *services.DonationService
	db     *gorm.DB
	config *viper.Viper
}

// CreateDonationRequest structure of the create request for donation
type CreateDonationRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
}

// NewDonationRouter creates a new instance of the donation router
func NewDonationRouter(app *gin.Engine) {
	r := DonationRouter{
		s:      services.NewDonationService(),
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}

	app.GET("/api/v1/donations", r.getAllDonations)
	app.GET("/api/v1/donations/:id", r.getDonationByID)
	app.POST("/api/v1/donations", r.createDonation)
	app.DELETE("/api/v1/donations/:id", r.deleteInvitation)
}

func (r *DonationRouter) getAllDonations(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, r.s.GetAll())
}

func (r *DonationRouter) createDonation(ctx *gin.Context) {
	var data CreateDonationRequest
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "missing or incorrect fields received", "error": "DONATION_CREATE_PAYLOAD_INVALID", "details": err.Error()})
		return
	}
	d := r.s.Create(data.FirstName, data.FirstName, data.Message, data.Amount)
	ctx.JSON(http.StatusOK, d)
}

func (r *DonationRouter) getDonationByID(ctx *gin.Context) {
	id := readID(ctx)
	if id != nil {
		d, err := r.s.GetByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "donation not found", "error": "DONATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, d)
	}
}

func (r *DonationRouter) deleteInvitation(ctx *gin.Context) {
	id := readID(ctx)
	if id != nil {
		i, err := r.s.DeleteByID(*id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "invitation not found", "error": "DONATION_NOT_FOUND", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, i)
	}
}
