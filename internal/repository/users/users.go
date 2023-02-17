package users

import (
	"context"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/conn"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RepositoryUsr struct {
	pool *pgxpool.Pool
}

func NewRepository(cfg config.Config) (*RepositoryUsr, error) {

	pool := conn.NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()
	if _, err := pool.Exec(ctx, `	
	create table if not exists users (	    
		user_login text unique not null,
		user_pass text not null,
		user_id varchar(27) unique not null
	);	
`); err != nil {
		return nil, err
	}

	return &RepositoryUsr{
		pool: pool,
	}, nil
}

func (p *RepositoryUsr) AddUserAuthData(ctx context.Context, login, pass, ID string) error {

	_, err := p.pool.Exec(ctx, `insert into users (user_login, user_pass, user_id) values ($1, $2, $3)`, login, pass, ID)
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				return model.ErrConflict
			}
		}

		return err
	}

	return nil
}

func (p *RepositoryUsr) GetUserAuthData(ctx context.Context, login, pass string) (string, error) {

	l, ps, ID := "", "", ""
	err := p.pool.QueryRow(ctx, `select user_login, user_pass, user_id from users where user_login = $1 and user_pass = $2`, login, pass).
		Scan(&l, &ps, &ID)
	if err != nil {
		return "", model.ErrNotExist
	}

	return ID, nil
}

func (p *RepositoryUsr) PingDB() error {

	pool := p.pool
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pool.Ping(ctx)

	return err
}
