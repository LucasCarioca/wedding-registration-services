package models

type Guest struct {
	Base
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	Email        string     `json:"email" binding:"required"`
	Phone        string     `json:"phone" binding:"required"`
	InvitationID uint       `json:"invitation_id" binding:"required"`
	Invitation   Invitation `json:"invitation" binding:"required"`
}
