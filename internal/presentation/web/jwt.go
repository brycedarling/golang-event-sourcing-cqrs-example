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

// CustomClaims ...
type CustomClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

const twentyFourHours time.Duration = time.Duration(60*60*24) * time.Second

// SignJWT ...
func SignJWT(userID string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(twentyFourHours).Unix(),
			Issuer:    "micro",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseJWT ...
func ParseJWT(authorization string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		authorization,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("JWT is expired")
	}
	return claims, nil
}
