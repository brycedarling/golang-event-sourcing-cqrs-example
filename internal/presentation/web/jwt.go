package web

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey []byte

func init() {
	secretKeyStr := os.Getenv("JWT_SECRET_KEY")
	if secretKeyStr == "" {
		secretKeyStr = "superSecureSecretText"
	}
	secretKey = []byte(secretKeyStr)
}

type customClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

const twentyFourHours time.Duration = time.Duration(60*60*24) * time.Second

func signJWT(userID string) (string, error) {
	claims := customClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(twentyFourHours).Unix(),
			Issuer:    "micro",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func parseJWT(authorization string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(
		authorization,
		&customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
