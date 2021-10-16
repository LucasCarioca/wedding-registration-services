package services

import (
	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_invitation_services(t *testing.T) {
	config.Init("dev")
	datasource.Init(config.GetConfig())
	is := NewInvitationService()

	t.Run("should create an invitation", func(t *testing.T) {
		invitationName := "test"
		invitationGuestCount := 1
		i := is.CreateInvitation(invitationName, invitationGuestCount)
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should get an invitation by id", func(t *testing.T) {
		invitationName := "test"
		invitationGuestCount := 1
		id := is.CreateInvitation(invitationName, invitationGuestCount).ID
		i, err := is.GetInvitationByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should get an invitation by registration key", func(t *testing.T) {
		invitationName := "test"
		invitationGuestCount := 1
		rk := is.CreateInvitation(invitationName, invitationGuestCount).RegistrationKey
		i, err := is.GetInvitationByRegistrationKey(rk)
		assert.Nilf(t, err, "should not throw an error")
		assert.NotNilf(t, i.ID, "should get an id")
		assert.NotNilf(t, i.Registered, "should generate a registration key")
		assert.Equal(t, invitationName, i.Name, "should have have the right name")
		assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
		is.DeleteInvitationByID(int(i.ID))
	})

	t.Run("should delete an invitation by id", func(t *testing.T) {
		invitationName := "test"
		invitationGuestCount := 1
		id := is.CreateInvitation(invitationName, invitationGuestCount).ID
		i, err := is.GetInvitationByID(int(id))
		assert.Nilf(t, err, "should not throw an error")
		is.DeleteInvitationByID(int(i.ID))
		_, err = is.GetInvitationByID(int(id))
		assert.NotNilf(t, err, "should throw an error")
		assert.Equal(t,"INVITATION_NOT_FOUND", err.Error(), "should throw the right error")
	})

	t.Run("should get all invitations", func(t *testing.T) {
		invitationName := "test"
		invitationGuestCount := 1
		is.CreateInvitation(invitationName, invitationGuestCount)
		is.CreateInvitation(invitationName, invitationGuestCount)
		is.CreateInvitation(invitationName, invitationGuestCount)
		invitations := is.GetAllInvitations()
		for _, i := range invitations {
			assert.NotNilf(t, i.Registered, "should generate a registration key")
			assert.Equal(t, invitationName, i.Name, "should have have the right name")
			assert.Equal(t, invitationGuestCount, i.GuestCount, "should have have the right guest count")
			is.DeleteInvitationByID(int(i.ID))
		}
	})
}