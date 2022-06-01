package services

import (
	"errors"

	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/LucasCarioca/wedding-registration-services/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//GuestService service for managing guests
type GuestService struct {
	db     *gorm.DB
	config *viper.Viper
}

//NewGuestService creates an instance of the guest service
func NewGuestService() *GuestService {
	return &GuestService{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}
}

//GetAll returns a list of all guests
func (s *GuestService) GetAll() []models.Guest {
	guests := make([]models.Guest, 0)
	s.db.Preload(clause.Associations).Find(&guests)
	return guests
}

//GetAllByInvitationID returns a list of all guests for a given invitation id
func (s *GuestService) GetAllByInvitationID(id uint) ([]models.Guest, error) {
	guests := make([]models.Guest, 0)
	s.db.Table("guests").Where("invitation_id = ?", id).Preload(clause.Associations).Find(&guests)
	return guests, nil
}

//GetCountByInvitationID returns a count of all guests for a given invitation id
func (s *GuestService) GetCountByInvitationID(id uint) int {
	var gc int64
	s.db.Table("guests").Where("invitation_id = ?", id).Count(&gc)
	return int(gc)
}

//Create creates a new guest and returns it
func (s *GuestService) Create(firstName string, lastName string, i models.Invitation) models.Guest {
	g := &models.Guest{
		FirstName:  firstName,
		LastName:   lastName,
		Invitation: i,
		RSVP:       false,
	}
	s.db.Create(g)
	return *g
}

//GetByID returns a guest by its id and returns it and an error if not found
func (s *GuestService) GetByID(id int) (*models.Guest, error) {
	g := models.Guest{}
	var c int64
	s.db.Preload(clause.Associations).Find(&g, id).Count(&c)
	if c > 0 {
		return &g, nil
	}
	return nil, errors.New("GUEST_NOT_FOUND")
}

//DeleteByID deletes a guest by its id and returns the deleted item and an error is it cannot be found
func (s *GuestService) DeleteByID(id int) (*models.Guest, error) {
	g := models.Guest{}
	var c int64
	s.db.Find(&g, id).Count(&c)
	if c < 1 {
		return nil, errors.New("GUEST_NOT_FOUND")
	}
	s.db.Delete(&g)
	return &g, nil
}
