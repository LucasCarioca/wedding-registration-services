package models

//Guest model for the guest table
type Guest struct {
	Base
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	InvitationID uint       `json:"invitation_id" binding:"required"`
	Invitation   Invitation `json:"invitation" binding:"required"`
	RSVP         bool       `json:"rsvp" binding:"required"`
}
