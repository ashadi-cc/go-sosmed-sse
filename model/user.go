package model

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type BaseUser struct {
	Email    string `json:"email" gorm:"unique_index" valid:"email,required"`
	Password string `json:"-" gorm:"type:text" valid:"required"`
}

type User struct {
	CommonModel
	BaseUser
	Posts []Post `gorm:"foreignkey:UserID" json:"-"`
}

type MyClaims struct {
	jwt.StandardClaims
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type Login struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}
