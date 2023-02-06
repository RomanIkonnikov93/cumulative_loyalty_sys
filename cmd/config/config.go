package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8081"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSecretKey         string `env:"JWT_SECRET_KEY" envDefault:"Secret-Key!"`
}

func GetConfig() (*Config, error) {

	cfg := &Config{}
	flag.StringVar(&cfg.RunAddress, "a", "127.0.0.1:8081", "RUN_ADDRESS")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "DATABASE_URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "ACCRUAL_SYSTEM_ADDRESS")
	flag.StringVar(&cfg.JWTSecretKey, "j", "", "JWT_SECRET_KEY")

	flag.Parse()
	err := env.Parse(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
