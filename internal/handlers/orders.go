package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/authjwt"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/validation"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
)

func PostOrdersHandler(rep repository.Pool, cfg config.Config, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "text/plain" {
			logger.Printf("%v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val, err := validation.OrderValid(string(b))
		if err != nil {
			logger.Printf("%v", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !val {
			logger.Printf("%v", http.StatusUnprocessableEntity)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := rep.Orders.GetUserIDbyOrder(r.Context(), string(b))
		if err != nil {
			if errors.Is(err, model.ErrNotExist) {
				err = rep.Orders.AddOrder(r.Context(), userID, string(b))
				if err != nil {
					logger.Error(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				return
			} else {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if user == userID {
			w.WriteHeader(http.StatusOK)
		} else {
			logger.Printf("%v", http.StatusConflict)
			w.WriteHeader(http.StatusConflict)
		}
	}
}

func GetOrdersHandler(rep repository.Pool, cfg config.Config, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		list, err := rep.Orders.GetOrdersByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, model.ErrNotExist) {
				logger.Printf("%v", http.StatusNoContent)
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			} else {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].UploadedAt.After(list[j].UploadedAt)
		})

		resp, err := json.Marshal(list)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
