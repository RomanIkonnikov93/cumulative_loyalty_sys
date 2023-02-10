package model

import "github.com/golang-jwt/jwt/v4"

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}
