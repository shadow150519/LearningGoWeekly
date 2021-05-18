package model

import "github.com/dgrijalva/jwt-go"

type CustomClaim struct {
	Username string `json:"username" binding:"required"`
	Age int `json:"age" binding:"required"`
	jwt.StandardClaims
}

