package repository

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrConflict = errors.New("conflict on insert")
	ErrNotExist = errors.New("not exist")
)

const TimeOut = time.Second * 5

type Pool struct {
	pool *pgxpool.Pool
}

func NewConnection(cfg config.Config) *pgxpool.Pool {

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
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
	    order_accrual text not null default '0.0',
		upload_time timestamp not null default current_timestamp
	);
	create table if not exists withdrawn (
	    user_id varchar(27) not null,
		order_id text not null,
		order_accrual text not null,
		processed_time timestamp not null default current_timestamp
	);
`); err != nil {
		return nil, err
	}

	return &Pool{
		pool: pool,
	}, nil
}

func (p *Pool) AddUserAuthData(ctx context.Context, login, pass, ID string) error {

	_, err := p.pool.Exec(ctx, `insert into users (user_login, user_pass, user_id) values ($1, $2, $3)`, login, pass, ID)
	if err != nil {
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

	_, err := p.pool.Exec(ctx, `insert into orders (user_id, order_id) values ($1, $2)`, userID, order)
	if err != nil {
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

func (p *Pool) GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error) {

	rows, err := p.pool.Query(ctx, `select user_id, order_id, order_status, order_accrual, upload_time from orders where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	list := make([]model.Order, rows.CommandTag().RowsAffected())

	for rows.Next() {

		var user, accrual string
		var order model.Order

		err = rows.Scan(&user, &order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		if accrual != "0.0" {
			num, err := strconv.ParseFloat(accrual, 64)
			if err != nil {
				return nil, err
			}
			order.Accrual = num
		}

		list = append(list, order)
	}

	if len(list) < 1 {
		return nil, ErrNotExist
	}

	return list, nil
}

func (p *Pool) GetOrdersForScanner() ([]model.Order, error) {

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()

	rows, err := p.pool.Query(ctx,
		`select order_id, order_status from orders where order_status=$1 or order_status=$2 or order_status=$3 `,
		"NEW", "REGISTERED", "PROCESSING")
	if err != nil {
		return nil, err
	}

	list := make([]model.Order, rows.CommandTag().RowsAffected())

	for rows.Next() {
		var status string
		var order model.Order

		err = rows.Scan(&order.Number, &status)
		if err != nil {
			return nil, err
		}

		list = append(list, order)
	}

	if len(list) < 1 {
		return nil, ErrNotExist
	}

	return list, nil
}

func (p *Pool) UpdateOrderData(ctx context.Context, status, accrual, order string) error {

	_, err := p.pool.Exec(ctx, `update orders set order_status = $1 , order_accrual = $2 where order_id = $3`, status, accrual, order)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) GetWithdrawnOrdersByUserID(ctx context.Context, userID string) ([]model.Withdrawn, error) {

	rows, err := p.pool.Query(ctx, `select user_id, order_id, order_accrual, processed_time from withdrawn where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	list := make([]model.Withdrawn, rows.CommandTag().RowsAffected())

	for rows.Next() {

		var user, accrual string
		var order model.Withdrawn

		err = rows.Scan(&user, &order.Order, &accrual, &order.ProcessedAt)
		if err != nil {
			return nil, err
		}

		if accrual != "0.0" {
			num, err := strconv.ParseFloat(accrual, 64)
			if err != nil {
				return nil, err
			}
			order.Accrual = num
		}

		list = append(list, order)
	}

	if len(list) < 1 {
		return nil, ErrNotExist
	}

	return list, nil
}

func (p *Pool) AddWithdrawnOrder(ctx context.Context, userID, order, sum string) error {

	_, err := p.pool.Exec(ctx, `insert into withdrawn (user_id, order_id, order_accrual) values ($1, $2, $3)`, userID, order, sum)
	if err != nil {
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

func (p *Pool) PingDB() error {

	pool := p.pool
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	err := pool.Ping(ctx)

	return err
}
