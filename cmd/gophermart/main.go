package main

import (
	"log"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/server"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/userrepository"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("GetConfig: ", err)
	}

	rep, err := userrepository.NewUserRepository(*cfg)
	if err != nil {
		log.Fatal("NewUserRepository: ", err)
	}

	err = server.StartServer(cfg, *rep)
	if err != nil {
		log.Fatal("StartServer: ", err)
	}
}
