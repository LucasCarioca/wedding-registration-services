package services

import (
	"errors"
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/LucasCarioca/wedding-registration-services/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//InvitationService service for managing invitations
type InvitationService struct {
	db     *gorm.DB
	config *viper.Viper
}

//NewInvitationService creates an instance of the invitation service
func NewInvitationService() *InvitationService {
	return &InvitationService{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}
}

//GetAllInvitations returns a list with all invitations
func (s *InvitationService) GetAllInvitations() []models.Invitation {
	invitations := make([]models.Invitation, 0)
	s.db.Find(&invitations)
	return invitations
}

//GetInvitationByRegistrationKey returns an invitation by its registration key and an error is it cannot be found
func (s *InvitationService) GetInvitationByRegistrationKey(key string) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Where("registration_key = ?", key).First(&i).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//GetInvitationByID returns an invitation by its id and an error is it cannot be found
func (s *InvitationService) GetInvitationByID(id int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//CreateInvitation creates a new invitation and returns it
func (s *InvitationService) CreateInvitation(name string, message string, email string, phone string, guestCount int) (*models.Invitation, error) {
	i := &models.Invitation{
		Name:       name,
		Message:    message,
		Phone:      phone,
		Email:      email,
		GuestCount: guestCount,
		Registered: false,
	}
	dbc := s.db.Create(i)
	if dbc.Error != nil {
		return nil, dbc.Error
	}
	return i, nil
}

//DeleteInvitationByID deletes an invitation by its id and returns the deleted item and an error is it cannot be found
func (s *InvitationService) DeleteInvitationByID(id int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	s.db.Delete(&i)
	return &i, nil
}
