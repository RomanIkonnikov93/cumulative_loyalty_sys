package repository

import (
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository/orders"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository/users"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository/withdrawn"
)

type Pool struct {
	Users     *users.Repository
	Orders    *orders.Repository
	Withdrawn *withdrawn.Repository
}

func NewReps(cfg config.Config) (*Pool, error) {

	u, err := users.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	o, err := orders.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	w, err := withdrawn.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Users:     u,
		Orders:    o,
		Withdrawn: w,
	}, nil
}
