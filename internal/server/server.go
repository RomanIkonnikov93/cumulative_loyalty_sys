package server

import (
	"log"

	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/handlers"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/userrepository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func StartServer(cfg *config.Config, rep userrepository.Pool) error {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.GzipRequest)
	r.Use(handlers.GzipResponse)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.RegisterHandler(rep, *cfg))
		r.Post("/login", handlers.LoginHandler(rep, *cfg))
		r.Group(func(r chi.Router) {
			r.Use(handlers.Auth(*cfg))
			r.Post("/orders", handlers.PostOrdersHandler())
			r.Get("/orders", handlers.GetOrdersHandler())
			r.Get("/balance", handlers.BalanceHandler())
			r.Post("/balance/withdraw", handlers.PostWithdrawHandler())
			r.Get("/withdrawals", handlers.GetWithdrawalsHandler())
			r.Get("/ping", handlers.PingDB(rep))
		})
	})

	log.Println("server running")
	err := http.ListenAndServe(cfg.RunAddress, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}
