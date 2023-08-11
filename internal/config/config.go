package config

import (
	"flag"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

type AddressConfig struct {
	Addr *string `env:"ADDRESS"`
}

type BaseConfig struct {
	AddressConfig

	ReportInterval *int `env:"REPORT_INTERVAL"`
	PollInterval   *int `env:"POLL_INTERVAL"`
}

func NewAddressConfig() AddressConfig {
	return AddressConfig{}
}

func NewBaseConfig() BaseConfig {
	return BaseConfig{}
}

func (cfg *AddressConfig) Parse() error {
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

func (cfg *BaseConfig) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.Addr == nil {
		cfg.Addr = new(string)
		flag.StringVar(cfg.Addr, "a", "http://localhost:8080", "http server addr")
	}

	if cfg.ReportInterval == nil {
		cfg.ReportInterval = new(int)
		flag.IntVar(cfg.ReportInterval, "r", 10, "report interval")
	}

	if cfg.PollInterval == nil {
		cfg.PollInterval = new(int)
		flag.IntVar(cfg.PollInterval, "p", 2, "poll interval")
	}

	flag.Parse()

	*cfg.Addr = normalizeAddr(*cfg.Addr)

	return nil
}

func normalizeAddr(addr string) string {
	if utf8.RuneCountInString(addr) < 4 || addr[:4] != "http" {
		return "http://" + addr
	}

	return addr
}
