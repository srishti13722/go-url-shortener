package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

//Generate jwt creates a jwt token for the use

func GenerateJWT(username string) (string, error){
	//define jwt claims payload
	claims := jwt.MapClaims{
		"username" : username,
		"exp" : time.Now().Add(time.Hour * 24).Unix(), //expiry 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ValidateJWT checks if a JWT token is valid
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}