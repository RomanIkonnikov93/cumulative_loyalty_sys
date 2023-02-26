package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/authjwt"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/model"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
	"github.com/segmentio/ksuid"
)

func RegisterHandler(rep repository.Pool, cfg config.Config, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
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

		data := model.UserAuth{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//generate new user ID
		ID := ksuid.New().String()

		err = rep.Users.AddUserAuthData(r.Context(), data.Login, data.Password, ID)
		if err != nil {
			if errors.Is(err, model.ErrConflict) {
				logger.Error(err)
				w.WriteHeader(http.StatusConflict)
				return
			} else {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		token, err := authjwt.EncodeJWT(ID, cfg.JWTSecretKey)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(rep repository.Pool, cfg config.Config, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "application/json" {
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

		data := model.UserAuth{}
		err = json.Unmarshal(b, &data)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ID, err := rep.Users.GetUserAuthData(r.Context(), data.Login, data.Password)
		if err != nil {
			if errors.Is(err, model.ErrNotExist) || errors.Is(err, model.ErrWrongPass) {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			} else {
				logger.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		token, err := authjwt.EncodeJWT(ID, cfg.JWTSecretKey)
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Authorization", token)
		w.WriteHeader(http.StatusOK)
	}
}
