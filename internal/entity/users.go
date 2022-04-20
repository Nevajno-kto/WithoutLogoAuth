package entity

const (
	Admin  = 1
	Client = 2
)

type User struct {
	Id           int    `json:"id"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
	Type         int    `json:"userType"`
}
