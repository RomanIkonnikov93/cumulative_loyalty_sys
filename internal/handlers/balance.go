package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/authjwt"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
)

func BalanceHandler(rep repository.Pool, cfg config.Config, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		list, err := rep.Orders.GetOrdersByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, model.ErrNotExist) {
				list = append(list, model.Order{
					Accrual: 0.0,
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		current := 0.0
		for _, order := range list {
			current += order.Accrual
		}

		wList, err := rep.Withdrawn.GetWithdrawnOrdersByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, model.ErrNotExist) {
				wList = append(wList, model.Withdrawn{
					Accrual: 0.0,
				})
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		withdrawn := 0.0
		for _, order := range wList {
			withdrawn += order.Accrual
		}

		data := model.Response{
			Current:   current - withdrawn,
			Withdrawn: withdrawn,
		}

		resp, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
