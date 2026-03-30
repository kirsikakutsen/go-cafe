package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getKey() string {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		panic("jwt secret not found")
	}
	return key
}

func SignPayLoad(accID uint) (string, error) {
	claims := jwt.MapClaims{
		"accID": accID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getKey()))
}

func VerifyPayload(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(getKey()), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithExpirationRequired(),
	)
}