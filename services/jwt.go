package services

import (
	"chat-service/models"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt methods
type JWTService interface {
	GenerateToken(user models.User) (string, string)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// data stored in the jwt token
type jwtCustomClaims struct {
	User models.User
	jwt.StandardClaims
}

// data stored in jwt refresh token (only id not all user info)
type jwtCustomClaimsRefreshToken struct {
	Id int
	jwt.StandardClaims
}

type jwtService struct {
	SecretKey string
	Issuer    string
}

func NewJwtService() JWTService {
	return &jwtService{
		SecretKey: getSecretKey(),
		Issuer:    "pragmaticreviews.com",
	}
}

// sets a default jwt secret if .env file JWT_SECRET key is empty for testing
func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = ";lkamfcacfp;m;cfma;lma;lm"
	}
	return secret
}

// token generation method
func (jwtSrv *jwtService) GenerateToken(user models.User) (string, string) {

	// both token and refresh token get created with their data
	claimsToken := &jwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtSrv.Issuer,
			Subject:   "token",
		},
	}
	claimsRefreshToken := &jwtCustomClaimsRefreshToken{
		Id: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 170).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtSrv.Issuer,
			Subject:   "refreshToken",
		},
	}

	// both tokens gets signed here with projects secret key before return
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsToken)
	secretToken, err := token.SignedString([]byte(jwtSrv.SecretKey))
	if err != nil {
		panic(err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)
	secretRefreshToken, err := refreshToken.SignedString([]byte(jwtSrv.SecretKey))
	if err != nil {
		panic(err)
	}
	return secretToken, secretRefreshToken
}

// token validation method used in the jwt authorization middleware
func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unxpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSrv.SecretKey), nil
	})

}
