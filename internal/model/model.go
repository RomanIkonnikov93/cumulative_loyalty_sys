package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type MyCustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
