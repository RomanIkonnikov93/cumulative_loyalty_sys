package withdrawn

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
	create table if not exists withdrawn (
	    user_id varchar(27) not null,
		order_id text not null,
		order_accrual text not null,
		processed_time timestamp not null default current_timestamp
	);
`); err != nil {
		return nil, err
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (p *Repository) GetWithdrawnOrdersByUserID(ctx context.Context, userID string) ([]model.Withdrawn, error) {

	rows, err := p.pool.Query(ctx, `select user_id, order_id, order_accrual, processed_time from withdrawn where user_id=$1`, userID)
	if err != nil {
		return nil, err
	}

	list := make([]model.Withdrawn, 0)

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
		return nil, model.ErrNotExist
	}

	return list, nil
}

func (p *Repository) AddWithdrawnOrder(ctx context.Context, userID, order, sum string) error {

	_, err := p.pool.Exec(ctx, `insert into withdrawn (user_id, order_id, order_accrual) values ($1, $2, $3)`, userID, order, sum)
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
