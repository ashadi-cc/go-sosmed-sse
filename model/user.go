package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type BaseUser struct {
	Email string `json:"email" gorm:"unique_index"`
	Password string `json:"password"`
}

type User struct {
	CommonModel
	BaseUser
}

type MyClaims struct {
	jwt.StandardClaims
	ID uint `json:"id"`
	Email string `json:"email"`
}
