package datasource

import (
	"fmt"
	"github.com/LucasCarioca/go-template/pkg/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

//Init initializes the database connection
func Init(config *viper.Viper) {
	var err error
	driver := config.GetString("data_source.driver")
	if driver == "sqlite" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	} else if driver == "postgres" {
		host := config.GetString("data_source.host")
		user := config.GetString("data_source.username")
		password := config.GetString("data_source.password")
		dbname := config.GetString("data_source.database")
		port := config.GetString("data_source.port")
		sslMode := config.GetString("data_source.ssl_mode")
		timeZone := config.GetString("data_source.time_zone")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslMode, timeZone)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
	if err != nil {
		panic("failed to connect database")
	}
	models.Init(db)
}

//GetDataSource gets the instance of the database connection
func GetDataSource() *gorm.DB {
	return db
}
