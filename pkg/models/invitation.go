package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invitation struct {
	gorm.Model
	Name            string `json:"Name" binding:"required"`
	GuestCount      int    `json:"guest_count" binding:"required"`
	Registered      bool   `json:"registered" binding:"required"`
	RegistrationKey string `json:"registration_key" binding:"required" gorm:"index;"`
}

func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	i.RegistrationKey = uuid.NewString()
	return nil
}