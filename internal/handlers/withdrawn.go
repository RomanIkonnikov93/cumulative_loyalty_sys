package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/authjwt"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/validation"
)

func PostWithdrawHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := model.WriteOff{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val, err := validation.OrderValid(data.Order)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if !val {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

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

		balance := current - withdrawn
		if balance < data.Sum {
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}

		log.Printf("PostWithdraw:%v,%v,%v", userID, data, balance)

		err = rep.Withdrawn.AddWithdrawnOrder(r.Context(), userID, data.Order, fmt.Sprintf("%g", data.Sum))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetWithdrawalsHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		list, err := rep.Withdrawn.GetWithdrawnOrdersByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, model.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].ProcessedAt.After(list[j].ProcessedAt)
		})

		resp, err := json.Marshal(list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("GetWithdrawals:%v,%v", userID, list)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
