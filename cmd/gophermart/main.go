package main

import (
	"context"
	"log"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/scanner"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/server"
)

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("GetConfig: ", err)
	}

	rep, err := repository.NewReps(*cfg)
	if err != nil {
		log.Fatal("NewReps: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err = scanner.Loop(ctx, *rep, *cfg)
		if err != nil {
			log.Fatal("scanner.Loop", err)
		}
	}()

	err = server.StartServer(*rep, *cfg)
	if err != nil {
		log.Fatal("StartServer: ", err)
	}
}
