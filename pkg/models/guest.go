package models

//Guest model for the guest table
type Guest struct {
	Base
	FirstName     string     `json:"first_name" binding:"required"`
	LastName      string     `json:"last_name" binding:"required"`
	Email         string     `json:"email" binding:"required"`
	Phone         string     `json:"phone" binding:"required"`
	EmailOptIn    bool       `json:"email_opt_in" binding:"required"`
	SMSOptIn      bool       `json:"sms_opt_in" binding:"required"`
	StreetAddress string     `json:"street_address" binding:"required"`
	City          string     `json:"city" binding:"required"`
	State         string     `json:"state" binding:"required"`
	ZipCode       string     `json:"zip_code" binding:"required"`
	InvitationID  uint       `json:"invitation_id" binding:"required"`
	Invitation    Invitation `json:"invitation" binding:"required"`
}
