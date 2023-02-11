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
	"github.com/segmentio/ksuid"
)

func RegisterHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
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

		data := model.UserAuth{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//generate new user ID
		ID := ksuid.New().String()

		err = rep.AddUserAuthData(r.Context(), data.Login, data.Password, ID)
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				w.WriteHeader(http.StatusConflict)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		token, err := authjwt.EncodeJWT(ID, cfg.JWTSecretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
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

		data := model.UserAuth{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ID, err := rep.GetUserAuthData(r.Context(), data.Login, data.Password)
		if err != nil {
			if errors.Is(err, repository.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		token, err := authjwt.EncodeJWT(ID, cfg.JWTSecretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)

	}
}

func PostOrdersHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !validation.OrderValid(string(b)) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := rep.GetUserIDbyOrder(r.Context(), string(b))
		if err != nil {
			if errors.Is(err, repository.ErrNotExist) {
				err = rep.AddOrder(r.Context(), userID, string(b))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusAccepted)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if user == userID {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	}
}

func GetOrdersHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		userID, err := authjwt.ParseJWTWithClaims(token, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		list, err := rep.GetOrdersByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, repository.ErrNotExist) {
				http.Error(w, err.Error(), http.StatusNoContent)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].UploadedAt.After(list[j].UploadedAt)
		})

		resp, err := json.Marshal(list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}

func BalanceHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func PostWithdrawHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetWithdrawalsHandler(rep repository.Pool, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func PingDB(rep repository.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rep.PingDb()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
