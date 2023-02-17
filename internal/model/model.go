package model

import (
	"errors"
	"time"
)

var (
	ErrConflict = errors.New("conflict on insert")
	ErrNotExist = errors.New("not exist")
	Err409      = errors.New("too many requests")
)

const TimeOut = time.Second * 10

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Withdrawn struct {
	Order       string    `json:"order,omitempty"`
	Accrual     float64   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

type Response struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WriteOff struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type ResponseForScanner struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
