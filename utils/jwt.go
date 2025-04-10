package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var SecretKey = []byte("your-secret-key")

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(SecretKey)
}
