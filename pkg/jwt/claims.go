package jwt

import (
	"fmt"

	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go/v4"
)

type Tokens struct {
	Refresh string `json:"refresh"`
	Auth    string `json:"auth"`
}

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

func SidnedTokens(authToken, refreshToken *jwt.Token, secret []byte) (Tokens, error) {
	var err error
	var tokens Tokens

	tokens.Auth, err = authToken.SignedString(secret)
	if err != nil {
		return tokens, errors.Wrap(entity.ErrServiceProblem, "usecase - auth - pkg - auth - SidnedTokens %w")
	}

	tokens.Refresh, err = authToken.SignedString(secret)
	if err != nil {
		return tokens, errors.Wrap(entity.ErrServiceProblem, "usecase - auth - pkg - refresg - SidnedTokens %w")
	}

	return tokens, err
}
