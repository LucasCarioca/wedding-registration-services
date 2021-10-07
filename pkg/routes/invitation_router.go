package routes

import (
	"github.com/LucasCarioca/go-template/pkg/datasource"
	"github.com/LucasCarioca/go-template/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type InvitationRouter struct {
	db *gorm.DB
}

type CreateInvitationRequest struct {
	Name       string `json:"name" binding:"required"`
	GuestCount int    `json:"guest_count" binding:"required"`
}

func NewInvitationRouter(app *gin.Engine) {
	r := InvitationRouter{
		db: datasource.GetDataSource(),
	}

	app.GET("/api/v1/invitations", r.getAllInvitations)
	app.GET("/api/v1/invitations/:id", r.getInvitation)
	app.POST("/api/v1/invitations", r.createInvitation)
	app.DELETE("/api/v1/invitations/:id", r.deleteInvitation)
}

func (r *InvitationRouter) readId(ctx *gin.Context) *int {
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
	invitations := make([]models.Invitation, 0)
	r.db.Find(&invitations)
	ctx.JSON(http.StatusOK, invitations)
}

func (r *InvitationRouter) createInvitation(ctx *gin.Context) {
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
	id := r.readId(ctx)
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
	id := r.readId(ctx)
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
