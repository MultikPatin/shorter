package models

import "github.com/golang-jwt/jwt/v4"

type User struct {
	ID       int64
	Username string
}

type Claims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}
