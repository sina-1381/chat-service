package models

import (
	"encoding/json"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgres connection data is created right after project boot on init
// func in main file and the connection data is stored in this global
// variable so it can be used all over the project . this way there is
// no need for a new db connection on every query
var DB *gorm.DB

// connecting to postgres based on .env info. this is called
// in init func in main file
func ConnectDataBase() {
	dsn := "host=" + os.Getenv("PGSQL_HOST") +
		" user=" + os.Getenv("PGSQL_USER") +
		" password=" + os.Getenv("PGSQL_PASS") +
		" sslmode=" + os.Getenv("PGSQL_SSLMODE") +
		" dbname=" + os.Getenv("PGSQL_DBNAME") +
		" TimeZone=" + os.Getenv("PGSQL_TIMEZONE")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // connection
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&User{}) // migration setup
	DB = database                 // connection info in DB global variable
}

// this converts different json models with similar keys
// mainly used this to convert database response data to final user response
func CreateResponse(model interface{}, response interface{}) {
	record, _ := json.Marshal(model)
	json.Unmarshal(record, response)
}
