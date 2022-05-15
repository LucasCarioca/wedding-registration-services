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
		i, _ := is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should get an invitation by id", func(t *testing.T) {
		inv, _ := is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		id := inv.ID
		i, err := is.GetInvitationByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should get an invitation by registration key", func(t *testing.T) {
		inv, _ := is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		rk := inv.RegistrationKey
		i, err := is.GetInvitationByRegistrationKey(rk)
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should delete an invitation by id", func(t *testing.T) {
		inv, _ := is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		id := inv.ID
		i, err := is.GetInvitationByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		is.DeleteInvitationByID(int(i.ID))
		_, err = is.GetInvitationByID(int(id))
		assert.NotNilf(t, err, "should throw an error")
		assert.Equal(t, "INVITATION_NOT_FOUND", err.Error(), "should throw the right error")
	})

	t.Run("should get all invitations", func(t *testing.T) {
		is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		is.CreateInvitation(invitationName, "testing message", email, phone, invitationGuestCount)
		invitations := is.GetAllInvitations()
		for _, i := range invitations {
			assert.NotNilf(t, i.Registered, "should generate a registration key")
			assert.Equal(t, invitationName, i.Name, "should have have the right name")
			assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
			is.DeleteInvitationByID(int(i.ID))
		}
	})
}
