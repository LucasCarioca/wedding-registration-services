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

// DonationService service for managing donations
type DonationService struct {
	db     *gorm.DB
	config *viper.Viper
}

// NewDonationService creates an instance of the donation service
func NewDonationService() *DonationService {
	return &DonationService{
		db:     datasource.GetDataSource(),
		config: config.GetConfig(),
	}
}

// GetAll returns a list of all guests
func (s *DonationService) GetAll() []models.Donation {
	donations := make([]models.Donation, 0)
	s.db.Preload(clause.Associations).Find(&donations)
	return donations
}

// Create creates a new guest and returns it
func (s *DonationService) Create(firstName string, lastName string, message string, amount string) models.Donation {
	d := &models.Donation{
		FirstName: firstName,
		LastName:  lastName,
		Message:   message,
		Amount:    amount,
	}
	s.db.Create(d)
	return *d
}

// GetByID returns a donation by its id and returns it and an error if not found
func (s *DonationService) GetByID(id int) (*models.Donation, error) {
	d := models.Donation{}
	var c int64
	s.db.Preload(clause.Associations).Find(&d, id).Count(&c)
	if c > 0 {
		return &d, nil
	}
	return nil, errors.New("DONATION_NOT_FOUND")
}

// DeleteByID deletes a donation by its id and returns the deleted item and an error is it cannot be found
func (s *DonationService) DeleteByID(id int) (*models.Donation, error) {
	d := models.Donation{}
	var c int64
	s.db.Find(&d, id).Count(&c)
	if c < 1 {
		return nil, errors.New("DONATION_NOT_FOUND")
	}
	s.db.Delete(&d)
	return &d, nil
}
