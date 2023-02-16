package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
)

var (
	Err409 = errors.New("too many requests")
)

func Loop(ctx context.Context, rep repository.Pool, cfg config.Config) error {

	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			_ = Scanner(rep, cfg)
		case <-ctx.Done():
			return nil
		}
	}
}

func Scanner(rep repository.Pool, cfg config.Config) error {

	list, err := rep.GetOrdersForScanner()
	if err != nil {
		return err
	}

	log.Printf("Scanner:%v", list)

	if len(list) < 1 {
		return nil
	}

	for _, order := range list {
		dur, err := UpdateOrders(rep, cfg, order.Number)
		if err != nil {
			if errors.Is(err, Err409) {
				time.Sleep(dur)
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func UpdateOrders(rep repository.Pool, cfg config.Config, order string) (time.Duration, error) {

	ctx, cansel := context.WithTimeout(context.Background(), repository.TimeOut)
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

	log.Printf("UpdateOrdersResp:%v", resp)

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

	log.Printf("UpdateOrderData:%v", data)

	err = rep.UpdateOrderData(ctx, data.Status, fmt.Sprintf("%g", data.Accrual), data.Order)
	if err != nil {
		return 0, nil
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		t := resp.Header.Get("Retry-After")
		dur, err := time.ParseDuration(t)
		if err != nil {
			dur = repository.TimeOut
		}
		return dur, Err409
	}

	return 0, nil
}
