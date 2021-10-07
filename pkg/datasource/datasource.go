package datasource

import (
	"github.com/LucasCarioca/go-template/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
func Init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	models.Init(db)
}

func GetDataSource() *gorm.DB {
	return db
}
