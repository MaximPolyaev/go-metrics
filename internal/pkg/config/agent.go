package config

import (
	"flag"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

type Agent struct {
	Addr           *string `env:"ADDRESS"`
	ReportInterval *int    `env:"REPORT_INTERVAL"`
	PollInterval   *int    `env:"POLL_INTERVAL"`
}

func NewAgent() Agent {
	return Agent{}
}

func (cfg *Agent) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	var isNeedParse bool

	if cfg.Addr == nil {
		flag.StringVar(cfg.Addr, "a", "http://localhost:8080", "http server addr")
		*cfg.Addr = normalizeAddr(*cfg.Addr)

		isNeedParse = true
	}

	if cfg.ReportInterval == nil {
		flag.IntVar(cfg.ReportInterval, "r", 10, "report interval")
		isNeedParse = true
	}

	if cfg.PollInterval == nil {
		flag.IntVar(cfg.PollInterval, "p", 2, "poll interval")
		isNeedParse = true
	}

	if isNeedParse {
		flag.Parse()
	}

	return nil
}

func normalizeAddr(addr string) string {
	if utf8.RuneCountInString(addr) < 4 || addr[:4] != "http" {
		return "http://" + addr
	}

	return addr
}
