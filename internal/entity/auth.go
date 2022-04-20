package entity

import "errors"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrServiceProblem     = errors.New("internal server error")
	ErrTimeout            = errors.New("request timeout")
)

const (
	SignUpRequest    = 1
	SignUpConfirm    = 2
	SignInRequest    = 3
	SignInConfirm    = 4
	SignInByPassword = 5
)

type IAuth interface{}

type Auth struct {
	Action int  `json:"action"`
	User   User `json:"user"`
	Code   int  `json:"code"`
}

type AuthTimeout struct {
	RemainingTime int64  `json:"remainingTime"`
	Msg           string `json:"message"`
}
