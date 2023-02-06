package userrepository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	AddUserAuthData(ctx context.Context, login, pass, token string) error
	GetUserAuthData(ctx context.Context, login, pass string) (string, error)
	PingDB(pool *pgxpool.Pool) error
}
