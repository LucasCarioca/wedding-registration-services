package services

import (
	"testing"

	"github.com/LucasCarioca/wedding-registration-services/pkg/config"
	"github.com/LucasCarioca/wedding-registration-services/pkg/datasource"
	"github.com/stretchr/testify/assert"
)

func Test_guest_services(t *testing.T) {
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
	emailOptIn := true
	smsOptIn := true
	street := "123 somewhere"
	city := "some-city"
	state := "some-state"
	zip := "55555"
	t.Run("should create a guest", func(t *testing.T) {
		i := is.CreateInvitation("test", 1)
		g := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		assert.NotNil(t, g.ID)
		assert.Equalf(t, i.ID, g.InvitationID, "should have a foreign key to the invitation")
		assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
		assert.Equalf(t, email, g.Email, "should have the right email")
		assert.Equalf(t, phone, g.Phone, "should have the right phone")
		assert.Equalf(t, emailOptIn, g.EmailOptIn, "should have the right email opt in")
		assert.Equalf(t, emailOptIn, g.SMSOptIn, "should have the right sms opt in")
		gs.DeleteGuestByID(int(g.ID))
	})

	t.Run("should search for a guest by id", func(t *testing.T) {
		i := is.CreateInvitation("test2", 1)
		id := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i).ID
		g, err := gs.GetGuestByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		assert.Equalf(t, i.ID, g.InvitationID, "should have a foreign key to the invitation")
		assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
		assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
		assert.Equalf(t, email, g.Email, "should have the right email")
		assert.Equalf(t, phone, g.Phone, "should have the right phone")
		gs.DeleteGuestByID(int(g.ID))
	})

	t.Run("should count the guests for an invitation", func(t *testing.T) {
		i := is.CreateInvitation("test3", 3)
		g1 := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		g2 := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		g3 := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		c := gs.GetGuestCountByInvitationID(i.ID)
		assert.Equalf(t, 3, c, "should show the right guest count number")
		gs.DeleteGuestByID(int(g1.ID))
		gs.DeleteGuestByID(int(g2.ID))
		gs.DeleteGuestByID(int(g3.ID))
	})

	t.Run("should be able to delete a guest by id", func(t *testing.T) {
		i := is.CreateInvitation("test4", 1)
		id := gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i).ID
		g, err := gs.GetGuestByID(int(id))
		assert.Nil(t, err, "should not throw an error")
		gs.DeleteGuestByID(int(g.ID))
		_, err = gs.GetGuestByID(int(id))
		assert.NotNil(t, err, "should throw an error")
		assert.Equalf(t, "GUEST_NOT_FOUND", err.Error(), "should throw the right error message")
	})

	t.Run("should get all guests for an invitation", func(t *testing.T) {
		i := is.CreateInvitation("test5", 3)
		gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i)
		guests, err := gs.GetAllGuestsByInvitationID(i.ID)
		assert.Nil(t, err, "should not throw an error")
		assert.Equalf(t, 3, len(guests), "should show the right number of guests in a list")
		for _, g := range guests {
			assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
			assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
			assert.Equalf(t, email, g.Email, "should have the right email")
			assert.Equalf(t, phone, g.Phone, "should have the right phone")
			gs.DeleteGuestByID(int(g.ID))
		}
	})

	t.Run("should get all guests", func(t *testing.T) {
		i1 := is.CreateInvitation("test6", 1)
		i2 := is.CreateInvitation("test6", 1)
		gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i1)
		gs.CreateGuest(firstName, lastName, email, phone, emailOptIn, smsOptIn, street, city, state, zip, i2)
		guests := gs.GetAllGuests()
		assert.Equalf(t, 2, len(guests), "should show the right number of guests in a list")
		for _, g := range guests {
			assert.Equalf(t, firstName, g.FirstName, "should have the right firstname")
			assert.Equalf(t, lastName, g.LastName, "should have the right lastname")
			assert.Equalf(t, email, g.Email, "should have the right email")
			assert.Equalf(t, phone, g.Phone, "should have the right phone")
			gs.DeleteGuestByID(int(g.ID))
		}
	})
}
