package utilities

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	// claims["exp"] = time.Now().Add(time.Second * 60).Unix()
	claims["iat"] = time.Now().Unix()

	return token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
}

func ParseAccessToken(token string) *jwt.Token {
	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	return parsedToken
}
