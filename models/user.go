package models

type User struct {
	Id          int64
	Email       string
	UserName    string
	Password    string
	IsActivated bool
	Role        string
}
