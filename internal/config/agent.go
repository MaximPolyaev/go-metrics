package config

import (
	"flag"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

// AgentConfig предназначен для настройки клиента сбора метрик
type AgentConfig struct {
	Addr           *string `env:"ADDRESS"`
	ReportInterval *int    `env:"REPORT_INTERVAL"`
	PollInterval   *int    `env:"POLL_INTERVAL"`
	HashKey        *string `env:"KEY"`
	RateLimit      *int    `env:"RATE_LIMIT"`
	CryptoKey      *string `env:"CRYPTO_KEY"`
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func (cfg *AgentConfig) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	addr := new(string)
	reportInterval := new(int)
	pollInterval := new(int)
	hashKey := new(string)
	rateLimit := new(int)
	cryptoKey := new(string)

	if cfg.Addr == nil {
		cfg.Addr = addr
	}
	if cfg.ReportInterval == nil {
		cfg.ReportInterval = reportInterval
	}
	if cfg.PollInterval == nil {
		cfg.PollInterval = pollInterval
	}
	if cfg.HashKey == nil {
		cfg.HashKey = hashKey
	}
	if cfg.RateLimit == nil {
		cfg.RateLimit = rateLimit
	}
	if cfg.CryptoKey == nil {
		cfg.CryptoKey = cryptoKey
	}

	flag.StringVar(addr, "a", "http://localhost:8080", "http server addr")
	flag.IntVar(reportInterval, "r", 10, "report interval")
	flag.IntVar(pollInterval, "p", 2, "poll interval")
	flag.StringVar(hashKey, "k", "", "hash key")
	flag.IntVar(rateLimit, "l", 0, "rate limit")
	flag.StringVar(cryptoKey, "crypto-key", "", "crypto key")

	flag.Parse()

	return nil
}

func (cfg *AgentConfig) GetNormalizedAddress() string {
	if nil == cfg.Addr {
		return ""
	}

	if utf8.RuneCountInString(*cfg.Addr) < 4 || (*cfg.Addr)[:4] != "http" {
		return "http://" + *cfg.Addr
	}

	return *cfg.Addr
}
