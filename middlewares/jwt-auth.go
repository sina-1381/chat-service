package middlewares

import (
	"ginGorm/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = authHeader[len(BEARER_SCHEMA):]
		} else if c.Param("token") != "" {
			tokenString = c.Param("token")
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		token, err := services.NewJwtService().ValidateToken(tokenString)
		claims := token.Claims.(jwt.MapClaims)
		if token.Valid && claims["sub"] == "token" {
			c.Set("user", claims)
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
