package orders

import (
	"context"
	"strconv"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/conn"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(cfg config.Config) (*Repository, error) {

	pool := conn.NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()
	if _, err := pool.Exec(ctx, `	
	create table if not exists orders (	    
		user_id varchar(27) not null,
		order_id text unique not null,
		order_status text not null default 'NEW',
	    order_accrual text not null default '0.0',
		upload_time timestamp not null default current_timestamp
	);	
`); err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (p *Repository) GetUserIDbyOrder(ctx context.Context, order string) (string, error) {

	user, ord := "", ""
	err := p.pool.QueryRow(ctx, `select user_id, order_id from orders where order_id = $1`, order).
		Scan(&user, &ord)
	if err != nil {
		return "", model.ErrNotExist
	}

	return user, nil
}

func (p *Repository) AddOrder(ctx context.Context, userID, order string) error {

	_, err := p.pool.Exec(ctx, `insert into orders (user_id, order_id) values ($1, $2)`, userID, order)
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

func (p *Repository) GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error) {

	rows, err := p.pool.Query(ctx, `select user_id, order_id, order_status, order_accrual, upload_time from orders where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	list := make([]model.Order, 0)

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
		return nil, model.ErrNotExist
	}

	return list, nil
}

func (p *Repository) GetOrdersForScanner() ([]model.Order, error) {

	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cancel()

	rows, err := p.pool.Query(ctx,
		`select order_id, order_status from orders where order_status=$1 or order_status=$2 or order_status=$3 `,
		"NEW", "REGISTERED", "PROCESSING")
	if err != nil {
		return nil, err
	}

	list := make([]model.Order, 0)

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
		return nil, model.ErrNotExist
	}

	return list, nil
}

func (p *Repository) UpdateOrderData(ctx context.Context, status, accrual, order string) error {

	_, err := p.pool.Exec(ctx, `update orders set order_status = $1 , order_accrual = $2 where order_id = $3`, status, accrual, order)
	if err != nil {
		return err
	}

	return nil
}
