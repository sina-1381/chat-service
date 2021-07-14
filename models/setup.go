package models

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDataBase() {
	dsn := "host="+os.Getenv("PGSQL_HOST")+" user="+os.Getenv("PGSQL_USER")+
		" password="+os.Getenv("PGSQL_PASS")+ " sslmode="+os.Getenv("PGSQL_SSLMODE")+
		" dbname="+os.Getenv("PGSQL_DBNAME")+" TimeZone="+os.Getenv("PGSQL_TIMEZONE")

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