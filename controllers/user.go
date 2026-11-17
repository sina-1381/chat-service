package controllers

import (
	"chat-service/models"
	"chat-service/services"
	"chat-service/validations"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// register api controller
func Register(c *gin.Context) {
	var model models.User
	var userResponse models.UserResponse

	// reciving request body and validating
	validations.CheckValidate(&model, c)

	// hashing the user input password before storing in db
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(model.Password), bcrypt.DefaultCost)
	model.Password = string(hash)

	models.DB.Create(&model) // data stored in db

	// creating api response out of db response (removing unwanted data)
	models.CreateResponse(model, &userResponse)

	// creating a 32 char random sting
	token := services.RandStringBytes(32)

	// caching key:user_email and value:random_string
	services.CacheRun(token, userResponse.Email)
	// this string is used in email verification making
	// sure the correct user is verifing their email

	// email send to the user
	services.TemplateEmailSender("/verify/email", token, userResponse.Email)

	// returning final json response
	c.JSON(http.StatusCreated, gin.H{
		"data":   userResponse,
		"status": "check your email to verify it",
	})
}

// verify email controller
func VerifyEmail(c *gin.Context) {
	var user models.User
	var email string

	// gets the token (random 32 char string we made earlier) from uri
	token := c.Param("token")

	// check the key value inside the cache to see if the token matches
	services.CacheRun(token, email)

	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "token invalid!",
		})
		return
	}

	// email is now verified in db
	r := models.DB.Model(user).Where("email=?", email).Update("status", true)

	r.First(&user) // get back user data from db

	// generate jwt based on user info
	jwtToken, jwtRefreshToken := services.NewJwtService().GenerateToken(user)

	// return final api response
	c.JSON(http.StatusOK,
		gin.H{"token": jwtToken, "refreshToken": jwtRefreshToken})
}

// user login controller
func Login(c *gin.Context) {
	var validateUser validations.Login
	var user models.User

	// input recive and validated
	validations.CheckValidate(&validateUser, c)

	// get user data based on email input
	r := models.DB.Model(user).Where("email=?", validateUser.Email).Scan(&user)
	if r.RowsAffected != 0 {

		// check account password
		err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(validateUser.Password))

		// error if password is invalid
		if err != nil {
			c.JSON(http.StatusForbidden, map[string]string{
				"error": "password invalid",
			})

		} else {

			// generate jwt tokens based on user info
			token, refreshToken := services.NewJwtService().GenerateToken(user)

			// jwt returned and user is now logged in
			c.JSON(http.StatusOK,
				gin.H{"token": token, "refreshToken": refreshToken})
		}
	} else {

		// if email does not exists response
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "this email does not exist",
		})
	}
}

// user check profile controller
func Profile(c *gin.Context) {
	var user models.User

	// get user id from their jwt token that they used in this request
	userId :=
		c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["ID"]

	// return user data found in db based on that user id
	models.DB.Model(user).Where("id=?", userId).Scan(&user)
	c.JSON(200, gin.H{"data": user})
}

// user change password controller
func ChangePassword(c *gin.Context) {
	var user models.User
	var password validations.ChangePassword

	// recive input and validate json body
	validations.CheckValidate(&password, c)

	// user id from jwt
	userId :=
		c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["ID"]

	// get user info from db with jwt id
	query := models.DB.Model(&user).Where("id = ?", userId).First(&user)

	// check if input old password matches the on in database
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(password.OldPassword))

	if err != nil {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "your old password is invalid",
		})
		return
	}

	// hash the new password from input
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(password.NewPassword), bcrypt.DefaultCost)

	// update the user password with new password in db
	query.Update("password", string(hash))

	// final response
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

// reset password api controller
func ResetPassword(c *gin.Context) {
	var input validations.ResetPassword

	// recive and validate request body
	validations.CheckValidate(&input, c)

	// generate random 32 char string
	token := services.RandStringBytes(32)

	// cache the string with users email as key
	services.CacheRun(token, input.Email)

	// send reset password email to user with the random string token attached
	services.TemplateEmailSender("/verify/password", token, input.Email)

	// final response
	c.JSON(http.StatusOK, map[string]string{
		"status": "check your email",
	})
}

// verify reset password api controller
func VerifyPassword(c *gin.Context) {
	var input validations.VerifyPassword
	var user models.User

	// recive request body and validate
	validations.CheckValidate(&input, c)

	var email string

	// check cache for user email key to get random string token
	services.CacheRun(input.Token, email)

	// if email and token are not cached or dont match return error
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "token invalid!",
		})
		return
	}

	// hash new input password
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(input.Password), bcrypt.DefaultCost)

	// update database with new hashed password
	models.DB.Model(&user).Where("email = ?", email).
		Update("password", string(hash))

	// final response
	c.JSON(http.StatusOK, gin.H{
		"message": "password changed.",
	})
}

// refresh token api controller
// this gets refresh token in input and geneates
// new token and refresh token based on that
func RefreshToken(c *gin.Context) {
	var user models.User
	const BEARER_SCHEMA = "Bearer "

	// get token from header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := authHeader[len(BEARER_SCHEMA):]

	// validate token
	token, err := services.NewJwtService().
		ValidateToken(tokenString)

	claims := token.Claims.(jwt.MapClaims)

	// if token is valid and verified :
	if token.Valid && claims["sub"] == "refreshToken" {

		// get user data from db based on user id inside refresh token
		models.DB.Model(&user).
			Where("id = ?", claims["Id"]).First(&user)

		// make new tokens based on user data inside db
		jwtToken, jwtRefreshToken :=
			services.NewJwtService().GenerateToken(user)

		// return tokens as final response
		c.JSON(http.StatusOK, gin.H{
			"token":        jwtToken,
			"refreshToken": jwtRefreshToken,
		})
	} else {

		// if token is not valid return error
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
