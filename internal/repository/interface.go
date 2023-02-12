package repository

import (
	"context"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddUserAuthData(ctx context.Context, login, pass, token string) error
	GetUserAuthData(ctx context.Context, login, pass string) (string, error)

	GetUserIDbyOrder(ctx context.Context, order string) (string, error)
	AddOrder(ctx context.Context, userID, order string) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error)
	GetWithdrawnOrdersByUserID(ctx context.Context, userID string) ([]model.Withdrawn, error)
	AddWithdrawnOrder(ctx context.Context, userID, order, sum string) error

	PingDB(pool *pgxpool.Pool) error
}
