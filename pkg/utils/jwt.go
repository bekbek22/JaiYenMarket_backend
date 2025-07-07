package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, role string, secret string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
