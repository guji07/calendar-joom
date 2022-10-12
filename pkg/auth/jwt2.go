package auth

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
)

type jwt2 struct {
	projectSecret string
}

func (j *jwt2) createToken(claims map[string]interface{}) (token string) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))

	token, err := newToken.SignedString(j.projectSecret)
	if err != nil {
		log.Fatal(err)
	}
	return token
}

func (j *jwt2) parseToken(token string) (err error) {
	return nil
}
