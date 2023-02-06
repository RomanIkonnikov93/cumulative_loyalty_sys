package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else {

				tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(cfg.JWTSecretKey), nil
				})
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if !tkn.Valid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					w.Header().Set("Authorization", token)
					next.ServeHTTP(w, r)
				}
			}
		})
	}
}

func GzipResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept-Encoding") == "gzip" {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

func GzipRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err := DecompressGZIP(b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(strings.NewReader(string(data)))
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
