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

//GetAll returns a list with all invitations
func (s *InvitationService) GetAll() []models.Invitation {
	invitations := make([]models.Invitation, 0)
	s.db.Find(&invitations)
	return invitations
}

//GetByRegistrationKey returns an invitation by its registration key and an error is it cannot be found
func (s *InvitationService) GetByRegistrationKey(key string) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Where("registration_key = ?", key).First(&i).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//GetByID returns an invitation by its id and an error is it cannot be found
func (s *InvitationService) GetByID(id int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//Create creates a new invitation and returns it
func (s *InvitationService) Create(name string, message string, email string, phone string, guestCount int) (*models.Invitation, error) {
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

//GetByPhone gets an invitation by phone
func (s *InvitationService) GetByPhone(phone string) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Where("phone = ?", phone).First(&i).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//GetByEmail gets an invitation by email
func (s *InvitationService) GetByEmail(email string) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Where("email = ?", email).First(&i).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	return &i, nil
}

//Search Attempts to find an invitation by phone email and registration key
func (s *InvitationService) Search(value string) (*models.Invitation, error) {
	i, err := s.GetByRegistrationKey(value)
	if err == nil {
		return i, nil
	}

	i, err = s.GetByPhone(value)
	if err == nil {
		return i, nil
	}

	i, err = s.GetByEmail(value)
	if err == nil {
		return i, nil
	}

	return nil, err
}

//DeleteByID deletes an invitation by its id and returns the deleted item or an error is it cannot be found
func (s *InvitationService) DeleteByID(id int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	s.db.Delete(&i)
	return &i, nil
}

//DeclineById decline an invitation by its id and returns the item or an error is it cannot be found
func (s *InvitationService) DeclineById(id int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	i.Declined = true
	s.db.Save(&i)
	return &i, nil
}

//UpdateGuestCountById update an invitation guest count by its id and returns the item or an error is it cannot be found
func (s *InvitationService) UpdateGuestCountById(id int, change int) (*models.Invitation, error) {
	i := models.Invitation{}
	var c int64
	s.db.Find(&i, id).Count(&c)
	if c < 1 {
		return nil, errors.New("INVITATION_NOT_FOUND")
	}
	i.GuestCount = i.GuestCount + change
	s.db.Save(&i)
	return &i, nil
}
