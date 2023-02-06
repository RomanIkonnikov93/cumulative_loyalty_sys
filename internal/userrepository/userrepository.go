package userrepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	ConnectTimeOut = time.Second * 5
	TimeOut        = time.Second * 3
)

var (
	ErrConflict = errors.New("conflict on insert")
)

type Pool struct {
	pool *pgxpool.Pool
}

func NewConnection(cfg config.Config) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), ConnectTimeOut)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return pool
}

func NewUserRepository(cfg config.Config) (*Pool, error) {

	pool := NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	if _, err := pool.Exec(ctx, `
	create table if not exists users (	    
	    user_login text unique,
	    user_pass text,
	    user_id varchar(27) unique
	)
`); err != nil {
		return nil, err
	}

	return &Pool{
		pool: pool,
	}, nil
}

func (p *Pool) AddUserAuthData(ctx context.Context, login, pass, ID string) error {

	if _, err := p.pool.Exec(ctx, `insert into users (user_login, user_pass, user_id) values ($1, $2, $3)`, login, pass, ID); err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		if ok {
			if pgerr.Code == "23505" {
				return ErrConflict
			}
		}

		return err
	}

	return nil
}

func (p *Pool) GetUserAuthData(ctx context.Context, login, pass string) (string, error) {

	rows, err := p.pool.Query(ctx, `select user_login, user_pass, user_id from users where user_login = $1 and user_pass = $2`, login, pass)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var l, ps, ID string

	for rows.Next() {
		if err := rows.Scan(&l, &ps, &ID); err != nil {
			return "", err
		}

	}

	return ID, nil
}

func (p *Pool) PingDb() error {
	pool := p.pool
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pool.Ping(ctx)
	return err
}
