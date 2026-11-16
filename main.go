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

func main() { // first layer of the project

	r := gin.Default() // gin framework setup

	r.SetTrustedProxies(nil) // trusted proxies setting part of the gin framework setup

	r.LoadHTMLGlob("templates/*") // loading up temp file for emails in templates folder

	ch := chat.NewSession() // setting up the websocket session for our live chat connections

	go ch.Run() // starting the websocket session in a seprate thread

	r.POST("/refresh/token", controllers.RefreshToken) // refresh token api route
	r.GET("/verify/email/:token", controllers.VerifyEmail) // verify email api route
	r.POST("/reset/password", controllers.ResetPassword) // reset password api route
	r.POST("/verify/password", controllers.VerifyPassword) // verify reset password api route
	r.POST("/login", controllers.Login) // login api route
	r.POST("/register", controllers.Register) // register api route

	userGroup := r.Group("/user") // "/user" prefix for these routes
	userGroup.Use(middlewares.AuthorizeJWT()) // these routes go under the jwt token verification middleware
	{
		userGroup.GET("/profile", controllers.Profile) // check profile api route
		userGroup.POST("/change/password", controllers.ChangePassword) // change password api route
	}

	wsGroup := r.Group("/ws") // "/ws" prefix for these routes
	wsGroup.Use(middlewares.AuthorizeJWT()) // these routes go under the jwt token verification middleware
	{
		wsGroup.GET("/:token", func(context *gin.Context) { // websocket live chat connection route
			chat.Runner(ch, context)
		})

		wsGroup.POST("/update/status", controllers.UpdateStatus) // update message status api route
		wsGroup.POST("/create/group", controllers.CreateGroup) // create group api route
		wsGroup.POST("/add/user", controllers.AddUser) // group add user api route
	}

	err := r.Run(os.Getenv("PORT")) // running the project on the port set in .env file
	if err != nil {
		panic(err)
	}
}

func init() { // this function runs first (before main) and only on project startup

	err := godotenv.Load(".env") // loading up .env file in the project
	if err != nil {
		panic(err)
	}
	
	gin.SetMode(os.Getenv("GIN_MODE")) // setting framework mode. debug or release

	// setting up mongodb connection with config set in env file
	err = mgm.SetDefaultConfig(nil, os.Getenv("MONGO_DBNAME"), options.Client().
		ApplyURI("mongodb://"+os.Getenv("MONGO_USER")+":"+os.Getenv("MONGO_PASS")+"@"+
			os.Getenv("MONGO_HOST")+":"+os.Getenv("MONGO_PORT")))
	if err != nil {
		panic(err)
	}

	models.ConnectDataBase() // calling the postgresql connection function
}
