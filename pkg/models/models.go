package models

import (
"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	db.AutoMigrate(&Invitation{})
	db.AutoMigrate(&Guest{})
}
