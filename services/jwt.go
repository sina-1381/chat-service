package services

import (
	"fmt"
	"ginGorm/models"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(user models.User) (string,string)
	ValidateToken(tokenString string)(*jwt.Token,error)
}
type jwtCustomClaims struct {
	User models.User
	jwt.StandardClaims
}

type jwtCustomClaimsRefreshToken struct {
	Id int
	jwt.StandardClaims
}

type jwtService struct {
	SecretKey string
	Issuer string
}

func NewJwtService() JWTService {
	return &jwtService{
		SecretKey: getSecretKey(),
		Issuer: "pragmaticreviews.com",
	}
}
func getSecretKey()  string {
	secret := os.Getenv("JWT_SECRET")
	if secret==""{
		secret = ";lkamfcacfp;m;cfma;lma;lm"
	}
	return secret
}
func (jwtSrv *jwtService) GenerateToken(user models.User)(string,string){
	claimsToken:=&jwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtSrv.Issuer,
			Subject: "token",
		},
	}
	claimsRefreshToken:= &jwtCustomClaimsRefreshToken{
		Id: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*24*170).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtSrv.Issuer,
			Subject: "refreshToken",
		},
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claimsToken)
	secretToken,err := token.SignedString([]byte(jwtSrv.SecretKey))
	if err!=nil{
		panic(err)
	}
	refreshToken:= jwt.NewWithClaims(jwt.SigningMethodHS256,claimsRefreshToken)
	secretRefreshToken,err := refreshToken.SignedString([]byte(jwtSrv.SecretKey))
	if err!=nil{
		panic(err)
	}
	return secretToken,secretRefreshToken
}
func (jwtSrv *jwtService)ValidateToken(tokenString string)(*jwt.Token,error){
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil,fmt.Errorf("unxpected signing method: %v",token.Header["alg"])
		}
		return []byte(jwtSrv.SecretKey),nil
	})

}