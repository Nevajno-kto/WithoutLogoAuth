package entity

type User struct {
	Id           int    `json:"id"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Organization string `json:"org"`
}
