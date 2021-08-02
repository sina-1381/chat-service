package main

import (
	"ginGorm/chat"
	"ginGorm/controllers"
	"ginGorm/middlewares"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	ch := chat.NewSession()
	go ch.Run()
	r.POST("/fcm", controllers.Fcm)
	r.POST("/refresh/token", controllers.RefreshToken)
	r.POST("/resend/email", controllers.ResendEmail)
	r.GET("/verify/email/:token", controllers.VerifyEmail)
	r.POST("/reset/password", controllers.ResetPassword)
	r.POST("/verify/password", controllers.VerifyPassword)
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)
	userGroup := r.Group("/user")
	userGroup.Use(middlewares.AuthorizeJWT())
	{
		userGroup.POST("/profile", controllers.Profile)
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
