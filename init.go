package main

import (
	"ginGorm/models"
	"ginGorm/validations"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func init()  {
	godotenv.Load(".env")
	gin.SetMode(os.Getenv("GIN_MODE"))

	err := mgm.SetDefaultConfig(nil , os.Getenv("MONGO_DBNAME") , options.Client().
		ApplyURI("mongodb://"+os.Getenv("MONGO_USER")+":"+os.Getenv("MONGO_PASS")+"@"+
			os.Getenv("MONGO_HOST")+":"+os.Getenv("MONGO_PORT")))

	if err != nil{
		panic(err)
	}
	models.ConnectDataBase()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("uniq", validations.Uniq)
		err = v.RegisterValidation("exists", validations.Exists)
		if err != nil{
			panic(err)
		}
	}
}
