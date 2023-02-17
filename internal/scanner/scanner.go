package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
)

func Loop(ctx context.Context, rep repository.Pool, cfg config.Config, logger logging.Logger) {

	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			scanner(rep, cfg, logger)
		case <-ctx.Done():
			logger.Info("Loop stopped")
			return
		}
	}
}

func scanner(rep repository.Pool, cfg config.Config, logger logging.Logger) {

	list, err := rep.Orders.GetOrdersForScanner()
	if err != nil {
		logger.Printf("scanner:%v", err)
	}

	if len(list) < 1 {
		return
	}

	for _, order := range list {
		dur, err := updateOrders(rep, cfg, order.Number)
		if err != nil {
			if errors.Is(err, model.Err409) {
				time.Sleep(dur)
				continue
			} else {
				return
			}
		}
	}
}

func updateOrders(rep repository.Pool, cfg config.Config, order string) (time.Duration, error) {

	ctx, cansel := context.WithTimeout(context.Background(), model.TimeOut)
	defer cansel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.AccrualSystemAddress+"/api/orders/"+order, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return 0, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	data := model.ResponseForScanner{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return 0, err
	}

	err = rep.Orders.UpdateOrderData(ctx, data.Status, fmt.Sprintf("%g", data.Accrual), data.Order)
	if err != nil {
		return 0, nil
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		t := resp.Header.Get("Retry-After")
		dur, err := time.ParseDuration(t)
		if err != nil {
			dur = model.TimeOut
		}
		return dur, model.Err409
	}

	return 0, nil
}
