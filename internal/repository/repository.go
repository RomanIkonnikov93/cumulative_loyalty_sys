package repository

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
	ErrNotExist = errors.New("not exist")
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

func NewRepository(cfg config.Config) (*Pool, error) {

	pool := NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	if _, err := pool.Exec(ctx, `
	create table if not exists users (	    
		user_login text unique not null,
		user_pass text not null,
		user_id varchar(27) unique not null
	);
	create table if not exists orders (	    
		user_id varchar(27) not null,
		order_id text unique not null,
		order_status text not null default 'NEW',
		upload_time timestamp not null default current_timestamp,
		order_accrual text not null default '0.0'
	);
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

	l, ps, ID := "", "", ""
	err := p.pool.QueryRow(ctx, `select user_login, user_pass, user_id from users where user_login = $1 and user_pass = $2`, login, pass).
		Scan(&l, &ps, &ID)
	if err != nil {
		return "", ErrNotExist
	}

	return ID, nil
}

func (p *Pool) GetUserIDbyOrder(ctx context.Context, order string) (string, error) {

	user, ord := "", ""
	err := p.pool.QueryRow(ctx, `select user_id, order_id from orders where order_id = $1`, order).
		Scan(&user, &ord)
	if err != nil {
		return "", ErrNotExist
	}

	return user, nil
}

func (p *Pool) AddOrder(ctx context.Context, userID, order string) error {

	if _, err := p.pool.Exec(ctx, `insert into orders (user_id, order_id) values ($1, $2)`, userID, order); err != nil {
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

func (p *Pool) PingDb() error {
	pool := p.pool
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pool.Ping(ctx)
	return err
}
