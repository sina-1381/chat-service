package main

import (
	"ginGorm/models"
	"ginGorm/validations"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init()  {
	err := mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 12} , "Message" , options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
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
