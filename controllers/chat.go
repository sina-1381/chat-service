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

// update message status api controller
// this api updates message status in mongodb.
// message status: send, recived, read
func UpdateStatus(c *gin.Context) {
	var input models.MessageStatus

	// validate input
	validations.CheckValidate(&input, c)

	// check if its a privet message or group message.
	// they are stored in different collections
	if input.Title == "" {
		go func() {

			// if privet message update status of the message
			// in message collection with input status
			username :=
				c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]
			_, err := mgm.Coll(&mongo.Message{}).UpdateMany(nil, bson.D{{"to", username},
				{"from", input.From}}, bson.D{{"$set", bson.D{{"status", input.Status}}}})
			if err != nil {
				panic(err)
			}
		}()
	} else {
		go func() {

			// if group message update status in group message collection with input status
			_, err := mgm.Coll(&mongo.GroupMessage{}).
				UpdateMany(nil, bson.D{{"title", input.Title}, {"from", input.From}},
					bson.D{{"$set", bson.D{{"status", input.Status}}}})
			if err != nil {
				panic(err)
			}
		}()
	}
}

// create group chat api controller
func CreateGroup(c *gin.Context) {
	var input models.GroupChat
	var user models.User

	// validate and recive input
	validations.CheckValidate(&input, c)

	// get user id from jwt token
	username :=
		c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["username"]

	// create group chat in database
	models.DB.Create(&input)

	// get user data from db
	models.DB.Table("users").Where("user_name = ?", username).First(&user)

	// create many2many connection between user and group chat in user_groups table
	data := models.UserGroups{
		GroupChatID: input.ID,
		UserID:      user.ID,
	}
	models.DB.Table("user_groups").Create(&data)

	// return final response
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

// add user to group chat api controller
func AddUser(c *gin.Context) {
	var input models.AddUser
	var group models.GroupChat

	// recive and validate request body
	validations.CheckValidate(&input, c)

	// find the input group in database
	models.DB.Table("group_chats").Where("title = ?", input.GroupTitle).First(&group)

	// add the input user id to the group
	data := models.UserGroups{
		GroupChatID: group.ID,
		UserID:      input.UserIDs,
	}
	models.DB.Table("user_groups").Create(&data)

	// return final response
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}
