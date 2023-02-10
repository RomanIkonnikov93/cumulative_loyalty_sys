package main

import (
	"log"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/server"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("GetConfig: ", err)
	}

	rep, err := repository.NewRepository(*cfg)
	if err != nil {
		log.Fatal("NewRepository: ", err)
	}

	err = server.StartServer(cfg, *rep)
	if err != nil {
		log.Fatal("StartServer: ", err)
	}
}
