package entity

import "errors"

var (
	ErrSignUp = errors.New("ошибка регистрации")
	ErrSingIn = errors.New("ошибка авторизации")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

type Auth struct {
	Type string `json:"type"`
	User User   `json:"user"`
	Code int    `json:"code"`
}
