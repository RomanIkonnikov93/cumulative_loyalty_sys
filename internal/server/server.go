package server

import (
	"log"

	"net/http"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/handlers"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func StartServer(rep repository.Pool, cfg config.Config, logger logging.Logger) error {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.GzipRequest)
	r.Use(handlers.GzipResponse)

	r.Route("/api/user", func(r chi.Router) {

		r.Post("/register", handlers.RegisterHandler(rep, cfg, logger))
		r.Post("/login", handlers.LoginHandler(rep, cfg, logger))

		r.Group(func(r chi.Router) {
			r.Use(handlers.Auth(cfg))

			r.Post("/orders", handlers.PostOrdersHandler(rep, cfg, logger))
			r.Get("/orders", handlers.GetOrdersHandler(rep, cfg, logger))
			r.Get("/balance", handlers.BalanceHandler(rep, cfg, logger))
			r.Post("/balance/withdraw", handlers.PostWithdrawHandler(rep, cfg, logger))
			r.Get("/withdrawals", handlers.GetWithdrawalsHandler(rep, cfg, logger))
			r.Get("/ping", handlers.PingDataBase(rep, logger))
		})
	})

	logger.Info("server running")
	err := http.ListenAndServe(cfg.RunAddress, r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}
