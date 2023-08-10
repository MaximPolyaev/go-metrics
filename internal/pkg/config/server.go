package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

type Server struct {
	Addr *string `env:"ADDRESS"`
}

func NewServer() Server {
	return Server{}
}

func (cfg *Server) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.Addr == nil {
		cfg.Addr = new(string)
		flag.StringVar(cfg.Addr, "a", ":8080", "http server addr")
		flag.Parse()
	}

	return nil
}
