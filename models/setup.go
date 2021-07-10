package models

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {
	dsn := "host=localhost user=postgres password=sina sslmode=disable dbname=postgres TimeZone=Asia/Tehran"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	database.AutoMigrate(&User{})
	DB = database
}

func CreateResponse(model interface{},response interface{}){
	record,_:=json.Marshal(model)
	json.Unmarshal(record, response)
}