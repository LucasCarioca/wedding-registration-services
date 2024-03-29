package services

import (
	"testing"

	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_invitation_services(t *testing.T) {
	config.Init("dev")
	datasource.Init(config.GetConfig())
	is := NewInvitationService()
	invitationName := "test"
	invitationGuestCount := 1
	email := "tester@email.com"
	phone := "5555555555"
	t.Run("should create an invitation", func(t *testing.T) {
		i, _ := is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteByID(int(i.ID))
	})

	t.Run("should get an invitation by id", func(t *testing.T) {
		inv, _ := is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		id := inv.ID
		i, err := is.GetByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteByID(int(i.ID))
	})

	t.Run("should get an invitation by registration key", func(t *testing.T) {
		inv, _ := is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		rk := inv.RegistrationKey
		i, err := is.GetByRegistrationKey(rk)
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteByID(int(i.ID))
	})

	t.Run("should delete an invitation by id", func(t *testing.T) {
		inv, _ := is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		id := inv.ID
		i, err := is.GetByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		is.DeleteByID(int(i.ID))
		_, err = is.GetByID(int(id))
		assert.NotNilf(t, err, "should throw an error")
		assert.Equal(t, "INVITATION_NOT_FOUND", err.Error(), "should throw the right error")
	})

	t.Run("should decline an invitation by id", func(t *testing.T) {
		inv, _ := is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		id := inv.ID
		i, _ := is.GetByID(int(id))
		assert.Equal(t, false, i.Declined, "should not have been declined")
		is.DeclineById(int(i.ID))
		i, _ = is.GetByID(int(id))
		assert.Equal(t, true, i.Declined, "should have been declined")
	})

	t.Run("should get all invitations", func(t *testing.T) {
		is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		is.Create(invitationName, "testing message", email, phone, invitationGuestCount)
		invitations := is.GetAll()
		for _, i := range invitations {
			assert.NotNilf(t, i.Registered, "should generate a registration key")
			assert.Equal(t, invitationName, i.Name, "should have the right name")
			assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
			is.DeleteByID(int(i.ID))
		}
	})
}
