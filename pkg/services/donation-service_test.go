package services

import (
	"testing"

	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_donation_services(t *testing.T) {
	config.Init("dev")
	datasource.Init(config.GetConfig())
	ds := NewDonationService()
	firstName := "tester"
	lastName := "tester"
	message := "somerandommessage"
	amount := "5555555555"
	t.Run("should create a donation", func(t *testing.T) {
		d := ds.CreateDonation(firstName, lastName, message, amount)
		assert.NotNil(t, d.ID)
		assert.Equalf(t, firstName, d.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, d.LastName, "should have the right lastname")
		assert.Equalf(t, lastName, d.Message, "should have the right message")
		assert.Equalf(t, lastName, d.Amount, "should have the right amount")
		ds.DeleteDonationByID(int(d.ID))
	})

	t.Run("should search for a donation by id", func(t *testing.T) {
		id := ds.CreateDonation(firstName, lastName, message, amount).ID
		d, err := ds.GetDonationByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		assert.Equalf(t, firstName, d.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, d.LastName, "should have the right lastname")
		ds.DeleteDonationByID(int(id))
	})

	t.Run("should be able to delete a donation by id", func(t *testing.T) {
		id := ds.CreateDonation(firstName, lastName, message, amount).ID
		d, err := ds.GetDonationByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		ds.DeleteDonationByID(int(d.ID))
		_, err = ds.GetDonationByID(int(id))
		assert.NotNil(t, err, "should throw an error")
		assert.Equalf(t, "DONATION_NOT_FOUND", err.Error(), "should throw the right error message")
	})


	t.Run("should get all donations", func(t *testing.T) {
		ds.CreateDonation(firstName, lastName, message, amount)
		ds.CreateDonation(firstName, lastName, message, amount)
		donations := ds.GetAllDonations()
		assert.Equalf(t, 2, len(donations), "should show the right number of donations in a list")
		for _, d := range donations {
			assert.Equalf(t, firstName, d.FirstName, "should have the right firstname")
			assert.Equalf(t, lastName, d.LastName, "should have the right lastname")
			ds.DeleteDonationByID(int(d.ID))
		}
	})
}
