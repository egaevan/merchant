package model

import jwt "github.com/dgrijalva/jwt-go"

type Token struct {
	UserID int
	Name   string
	Email  string
	*jwt.StandardClaims
}
