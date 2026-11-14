package main

import (
	"chat-service/chat"
	"chat-service/controllers"
	"chat-service/middlewares"
	"chat-service/models"
	"os"

	"github.com/Kamva/mgm"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*")
	ch := chat.NewSession()
	go ch.Run()
	r.POST("/refresh/token", controllers.RefreshToken)
	r.GET("/verify/email/:token", controllers.VerifyEmail)
	r.POST("/reset/password", controllers.ResetPassword)
	r.POST("/verify/password", controllers.VerifyPassword)
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	userGroup := r.Group("/user")
	userGroup.Use(middlewares.AuthorizeJWT())
	{
		userGroup.GET("/profile", controllers.Profile)
		userGroup.POST("/change/password", controllers.ChangePassword)
	}
	wsGroup := r.Group("/ws")
	wsGroup.Use(middlewares.AuthorizeJWT())
	{
		wsGroup.GET("/:token", func(context *gin.Context) {
			chat.Runner(ch, context)
		})
		wsGroup.POST("/update/status", controllers.UpdateStatus)
		wsGroup.POST("/create/group", controllers.CreateGroup)
		wsGroup.POST("/add/user", controllers.AddUser)
	}
	err := r.Run(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	err = mgm.SetDefaultConfig(nil, os.Getenv("MONGO_DBNAME"), options.Client().
		ApplyURI("mongodb://"+os.Getenv("MONGO_USER")+":"+os.Getenv("MONGO_PASS")+"@"+
			os.Getenv("MONGO_HOST")+":"+os.Getenv("MONGO_PORT")))
	if err != nil {
		panic(err)
	}
	models.ConnectDataBase()
}
