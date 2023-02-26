package orders

import (
	"context"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
)

type Orders interface {
	GetUserIDbyOrder(ctx context.Context, order string) (string, error)
	AddOrder(ctx context.Context, userID, order string) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error)
	GetOrdersForScanner() ([]model.Order, error)
	UpdateOrderData(ctx context.Context, status, accrual, order string) error
}
