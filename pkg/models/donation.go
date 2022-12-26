package models

// Donation model for the donation table
type Donation struct {
	Base
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
}
