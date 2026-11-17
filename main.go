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

	// trusted proxies setting part of the gin framework setup
	r.SetTrustedProxies(nil)

	// loading up temp file for emails in templates folder
	r.LoadHTMLGlob("templates/*")

	// setting up the websocket session for our live chat connections
	ch := chat.NewSession()

	go ch.Run() // starting the websocket session in a seprate thread

	// refresh token api route
	r.POST("/refresh/token", controllers.RefreshToken)

	// verify email api route
	r.GET("/verify/email/:token", controllers.VerifyEmail)

	// reset password api route
	r.POST("/reset/password", controllers.ResetPassword)

	// verify reset password api route
	r.POST("/verify/password", controllers.VerifyPassword)

	// login api route
	r.POST("/login", controllers.Login)

	// register api route
	r.POST("/register", controllers.Register)

	// "/user" prefix for these routes
	userGroup := r.Group("/user")

	// these routes go under the jwt token verification middleware
	userGroup.Use(middlewares.AuthorizeJWT())
	{
		// check profile api route
		userGroup.GET("/profile", controllers.Profile)

		// change password api route
		userGroup.POST("/change/password", controllers.ChangePassword)

		wsGroup := r.Group("/ws") // "/ws" prefix for these routes

		// these routes go under the jwt token verification middleware
		wsGroup.Use(middlewares.AuthorizeJWT())
		{
			// websocket live chat connection route
			wsGroup.GET("/:token", func(context *gin.Context) {
				chat.Runner(ch, context)
			})

			// update message status api route
			wsGroup.POST("/update/status", controllers.UpdateStatus)

			// create group api route
			wsGroup.POST("/create/group", controllers.CreateGroup)

			// group add user api route
			wsGroup.POST("/add/user", controllers.AddUser)
		}

		// running the project on the port set in .env file
		err := r.Run(os.Getenv("PORT"))
		if err != nil {
			panic(err)
		}
	}
}

// this function runs first (before main) and only on project startup
func init() {

	// loading up .env file in the project
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// setting framework mode. debug or release
	gin.SetMode(os.Getenv("GIN_MODE"))

	// setting up mongodb connection with config set in env file
	err = mgm.SetDefaultConfig(nil, os.Getenv("MONGO_DBNAME"), options.Client().
		ApplyURI("mongodb://"+os.Getenv("MONGO_USER")+":"+os.Getenv("MONGO_PASS")+
			"@"+os.Getenv("MONGO_HOST")+":"+os.Getenv("MONGO_PORT")))

	if err != nil {
		panic(err)
	}

	// calling the postgresql connection function
	models.ConnectDataBase()
}
