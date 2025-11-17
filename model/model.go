package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique" json:"username"`
	Email     string `gorm:"unique" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Password  string `json:"-" form:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email`
	Password string `json:"password`
}

type UserLoginResponse struct {
	Message  string `json:"message"`
	JWTToken string
}

type UserLogoutRequest struct {
	Email string `json:"email"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}
