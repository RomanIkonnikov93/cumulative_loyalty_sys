package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddUserAuthData(ctx context.Context, login, pass, token string) error
	GetUserAuthData(ctx context.Context, login, pass string) (string, error)

	GetUserIDbyOrder(ctx context.Context, order string) (string, error)
	AddOrder(ctx context.Context, userID, order string) error

	PingDB(pool *pgxpool.Pool) error
}
