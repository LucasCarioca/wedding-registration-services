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
	is := NewInvitationService()
	gs := NewGuestService()
	defer func() {
		invitations := is.GetAllInvitations()
		for _, i := range invitations {
			is.DeleteInvitationByID(int(i.ID))
		}
	}()
	firstName := "tester"
	lastName := "tester"
	email := "tester@email.com"
	phone := "5555555555"
	t.Run("should create a guest", func(t *testing.T) {
		i, _ := is.CreateInvitation("test", "testing message", email, phone, 1)
		g := gs.CreateGuest(firstName, lastName, *i)
		assert.NotNil(t, g.ID)
		assert.Equalf(t, i.ID, g.InvitationID, "should have a foreign key to the invitation")
		assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
		gs.DeleteGuestByID(int(g.ID))
	})

	t.Run("should search for a guest by id", func(t *testing.T) {
		i, _ := is.CreateInvitation("test2", "testing message", email, phone, 1)
		id := gs.CreateGuest(firstName, lastName, *i).ID
		g, err := gs.GetGuestByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		assert.Equalf(t, i.ID, g.InvitationID, "should have a foreign key to the invitation")
		assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
		gs.DeleteGuestByID(int(g.ID))
	})

	t.Run("should count the guests for an invitation", func(t *testing.T) {
		i, _ := is.CreateInvitation("test3", "testing message", email, phone, 3)
		g1 := gs.CreateGuest(firstName, lastName, *i)
		g2 := gs.CreateGuest(firstName, lastName, *i)
		g3 := gs.CreateGuest(firstName, lastName, *i)
		c := gs.GetGuestCountByInvitationID(i.ID)
		assert.Equalf(t, 3, c, "should show the right guest count number")
		gs.DeleteGuestByID(int(g1.ID))
		gs.DeleteGuestByID(int(g2.ID))
		gs.DeleteGuestByID(int(g3.ID))
	})

	t.Run("should be able to delete a guest by id", func(t *testing.T) {
		i, _ := is.CreateInvitation("test4", "testing message", email, phone, 1)
		id := gs.CreateGuest(firstName, lastName, *i).ID
		g, err := gs.GetGuestByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		gs.DeleteGuestByID(int(g.ID))
		_, err = gs.GetGuestByID(int(id))
		assert.NotNil(t, err, "should throw an error")
		assert.Equalf(t, "GUEST_NOT_FOUND", err.Error(), "should throw the right error message")
	})

	t.Run("should get all guests for an invitation", func(t *testing.T) {
		i, _ := is.CreateInvitation("test5", "testing message", email, phone, 3)
		gs.CreateGuest(firstName, lastName, *i)
		gs.CreateGuest(firstName, lastName, *i)
		gs.CreateGuest(firstName, lastName, *i)
		guests, err := gs.GetAllGuestsByInvitationID(i.ID)
		assert.Nil(t, err, "should not throw an error")
		assert.Equalf(t, 3, len(guests), "should show the right number of guests in a list")
		for _, g := range guests {
			assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
			assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
			gs.DeleteGuestByID(int(g.ID))
		}
	})

	t.Run("should get all guests", func(t *testing.T) {
		i1, _ := is.CreateInvitation("test6", "testing message", email, phone, 1)
		i2, _ := is.CreateInvitation("test6", "testing message", email, phone, 1)
		gs.CreateGuest(firstName, lastName, *i1)
		gs.CreateGuest(firstName, lastName, *i2)
		guests := gs.GetAllGuests()
		assert.Equalf(t, 2, len(guests), "should show the right number of guests in a list")
		for _, g := range guests {
			assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
			assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
			gs.DeleteGuestByID(int(g.ID))
		}
	})
}
