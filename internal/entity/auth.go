package entity

import "errors"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrServiceProblem = errors.New("internal server error")
)

type Auth struct {
	Type string `json:"type"`
	User User   `json:"user"`
	Code int    `json:"code"`
}
