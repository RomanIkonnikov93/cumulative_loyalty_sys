package withdrawn

import (
	"context"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
)

type Withdrawn interface {
	GetWithdrawnOrdersByUserID(ctx context.Context, userID string) ([]model.Withdrawn, error)
	AddWithdrawnOrder(ctx context.Context, userID, order, sum string) error
}
