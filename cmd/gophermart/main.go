package main

import (
	"context"

	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/cmd/config"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/repository"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/scanner"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/internal/server"
	"github.com/RomanIkonnikov93/cumulative_loyalty_sys/logging"
)

func main() {

	logger := logging.GetLogger()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Fatalf("GetConfig: %s", err)
	}

	rep, err := repository.NewReps(*cfg)
	if err != nil {
		logger.Fatalf("NewReps: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		scanner.Loop(ctx, *rep, *cfg, *logger)
	}()

	err = server.StartServer(*rep, *cfg, *logger)
	if err != nil {
		logger.Fatalf("StartServer: %s", err)
	}
}
