package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type BaseUser struct {
	Email    string `json:"email" gorm:"unique_index" valid:"email,required"`
	Password string `json:"password" gorm:"type:text" valid:"required"`
}

type User struct {
	CommonModel
	BaseUser
}

type MyClaims struct {
	jwt.StandardClaims
	ID    uint   `json:"id"`
	Email string `json:"email"`
}
