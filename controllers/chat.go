package controllers

import (
	"chat-service/models"
	"chat-service/mongo"
	"chat-service/validations"
	"net/http"

	"github.com/Kamva/mgm"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateStatus(c *gin.Context) {
	var input models.MessageStatus
	validations.CheckValidate(&input, c)
	if input.Title == "" {
		go func() {
			username := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]
			_, err := mgm.Coll(&mongo.Message{}).UpdateMany(nil, bson.D{{"to", username}, {"from", input.From}}, bson.D{{"$set", bson.D{{"status", input.Status}}}})
			if err != nil {
				panic(err)
			}
		}()
	} else {
		go func() {
			_, err := mgm.Coll(&mongo.GroupMessage{}).UpdateMany(nil, bson.D{{"title", input.Title}, {"from", input.From}},
				bson.D{{"$set", bson.D{{"status", input.Status}}}})
			if err != nil {
				panic(err)
			}
		}()
	}
}

func CreateGroup(c *gin.Context) {
	var input models.GroupChat
	var user models.User
	validations.CheckValidate(&input, c)
	username := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]
	models.DB.Create(&input)
	models.DB.Table("users").Where("user_name = ?", username).First(&user)
	data := models.UserGroups{
		GroupChatID: input.ID,
		UserID:      user.ID,
	}
	models.DB.Table("user_groups").Create(&data)
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func AddUser(c *gin.Context) {
	var input models.AddUser
	var group models.GroupChat
	validations.CheckValidate(&input, c)
	models.DB.Table("group_chats").Where("title = ?", input.GroupTitle).First(&group)
	data := models.UserGroups{
		GroupChatID: group.ID,
		UserID:      input.UserIDs,
	}
	models.DB.Table("user_groups").Create(&data)
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
