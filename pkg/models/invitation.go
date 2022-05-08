package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Invitation model for the invitation table
type Invitation struct {
	Base
	Name            string `json:"name" binding:"required"`
	Message         string `json:"message" binding:"required"'`
	GuestCount      int    `json:"guest_count" binding:"required"`
	Registered      bool   `json:"registered" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Email           string `json:"email" binding:"required"`
	RegistrationKey string `json:"registration_key" binding:"required" gorm:"index;"`
}

//BeforeCreate creates a random uuid registration key for new invitations
func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	i.RegistrationKey = uuid.NewString()
	return nil
}
