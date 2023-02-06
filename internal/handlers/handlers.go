package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/authjwt"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/userrepository"
	"github.com/segmentio/ksuid"
)

func RegisterHandler(rep userrepository.Pool, cfg config.Config) http.HandlerFunc {
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
			if errors.Is(err, userrepository.ErrConflict) {
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

func LoginHandler(rep userrepository.Pool, cfg config.Config) http.HandlerFunc {
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
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
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

func PostOrdersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetOrdersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func BalanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.BasicAuth()
	}
}

func PostWithdrawHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetWithdrawalsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func PingDB(rep userrepository.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rep.PingDb()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
