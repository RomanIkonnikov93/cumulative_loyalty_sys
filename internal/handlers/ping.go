package handlers

import (
	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
)

func PingDataBase(rep repository.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rep.Users.PingDB()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}
