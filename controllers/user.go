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

func Register(c *gin.Context) {
	var model models.User
	var userResponse models.UserResponse
	validations.CheckValidate(&model, c)
	hash, _ := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	model.Password = string(hash)
	models.DB.Create(&model)
	models.CreateResponse(model, &userResponse)
	token := services.RandStringBytes(32)
	services.CacheRun(token, userResponse.Email)
	services.TemplateEmailSender("/verify/email", token, userResponse.Email)
	c.JSON(http.StatusCreated, gin.H{
		"data":   userResponse,
		"status": "check your email to verify it",
	})
}

func VerifyEmail(c *gin.Context) {
	var user models.User
	var email string
	token := c.Param("token")
	services.CacheRun(token, email)
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "token invalid!",
		})
		return
	}
	r := models.DB.Model(user).Where("email=?", email).Update("status", true)
	r.First(&user)
	jwtToken, jwtRefreshToken := services.NewJwtService().GenerateToken(user)
	c.JSON(http.StatusOK, gin.H{"token": jwtToken, "refreshToken": jwtRefreshToken})
}

func Login(c *gin.Context) {
	var validateUser validations.Login
	var user models.User
	validations.CheckValidate(&validateUser, c)
	r := models.DB.Model(user).Where("email=?", validateUser.Email).Scan(&user)
	if r.RowsAffected != 0 {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(validateUser.Password))
		if err != nil {
			c.JSON(http.StatusForbidden, map[string]string{
				"error": "password invalid",
			})
		} else {
			token, refreshToken := services.NewJwtService().GenerateToken(user)
			c.JSON(http.StatusOK, gin.H{"token": token, "refreshToken": refreshToken})
		}
	} else {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "this email does not exist",
		})
	}
}

func Profile(c *gin.Context) {
	var user models.User
	userId := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["ID"]
	models.DB.Model(user).Where("id=?", userId).Scan(&user)
	c.JSON(200, gin.H{"data": user})
}

func ChangePassword(c *gin.Context) {
	var user models.User
	var password validations.ChangePassword
	validations.CheckValidate(&password, c)
	userId := c.MustGet("user").(jwt.MapClaims)["User"].(map[string]interface{})["ID"]
	query := models.DB.Model(&user).Where("id = ?", userId).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password.OldPassword))
	if err != nil {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "your old password is invalid",
		})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	query.Update("password", string(hash))
	c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func ResetPassword(c *gin.Context) {
	var input validations.ResetPassword
	validations.CheckValidate(&input, c)
	token := services.RandStringBytes(32)
	services.CacheRun(token, input.Email)
	services.TemplateEmailSender("/verify/password", token, input.Email)
	c.JSON(http.StatusOK, map[string]string{
		"status": "check your email",
	})
}

func VerifyPassword(c *gin.Context) {
	var input validations.VerifyPassword
	var user models.User
	validations.CheckValidate(&input, c)
	var email string
	services.CacheRun(input.Token, email)
	if email == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "token invalid!",
		})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	models.DB.Model(&user).Where("email = ?", email).Update("password", string(hash))
	c.JSON(http.StatusOK, gin.H{
		"message": "password changed.",
	})
}

func RefreshToken(c *gin.Context) {
	var user models.User
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	tokenString := authHeader[len(BEARER_SCHEMA):]
	token, err := services.NewJwtService().ValidateToken(tokenString)
	claims := token.Claims.(jwt.MapClaims)
	if token.Valid && claims["sub"] == "refreshToken" {
		models.DB.Model(&user).Where("id = ?", claims["Id"]).First(&user)
		jwtToken, jwtRefreshToken := services.NewJwtService().GenerateToken(user)
		c.JSON(http.StatusOK, gin.H{
			"token":        jwtToken,
			"refreshToken": jwtRefreshToken,
		})
	} else {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
