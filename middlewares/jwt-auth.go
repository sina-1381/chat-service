package middlewares

import (
	"chat-service/services"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// user authorization middleware with jwt token validation
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization") // recives the token from request header

		// handeling if token is not in header
		if authHeader != "" {
			tokenString = authHeader[len(BEARER_SCHEMA):]
		} else if c.Param("token") != "" {
			tokenString = c.Param("token")
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// validating jwt token with ValidateToken method in services/jwt file
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
