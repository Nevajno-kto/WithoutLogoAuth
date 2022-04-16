package jwt

import (
	"fmt"

	"github.com/nevajno-kto/without-logo-auth/internal/entity"

	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	jwt.StandardClaims
	UserId int `json:"UserId"`
}

func ParseToken(accessToken string, signingKey []byte) (string, error) {

	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.ID, nil
	}

	return "", entity.ErrInvalidAccessToken
}
