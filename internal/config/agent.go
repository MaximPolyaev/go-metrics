package config

import (
	"encoding/json"
	"flag"
	"os"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

const (
	defaultAgentAddr      = "http://localhost:8080"
	defaultAgentCryptoKey = ""
	defaultReportInterval = 10
	defaultPoolInterval   = 2
)

// AgentConfig предназначен для настройки клиента сбора метрик
type AgentConfig struct {
	Addr           *string `env:"ADDRESS"`
	ReportInterval *int    `env:"REPORT_INTERVAL"`
	PollInterval   *int    `env:"POLL_INTERVAL"`
	HashKey        *string `env:"KEY"`
	RateLimit      *int    `env:"RATE_LIMIT"`
	CryptoKey      *string `env:"CRYPTO_KEY"`
	JSONConfig     *string `env:"CONFIG"`
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
	jsonConfig := new(string)

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
	if cfg.JSONConfig == nil {
		cfg.JSONConfig = jsonConfig
	}

	flag.StringVar(addr, "a", defaultAgentAddr, "http server addr")
	flag.IntVar(reportInterval, "r", defaultReportInterval, "report interval")
	flag.IntVar(pollInterval, "p", defaultPoolInterval, "poll interval")
	flag.StringVar(hashKey, "k", "", "hash key")
	flag.IntVar(rateLimit, "l", 0, "rate limit")
	flag.StringVar(cryptoKey, "crypto-key", defaultAgentCryptoKey, "crypto key")
	flag.StringVar(jsonConfig, "c", "", "json config")
	flag.StringVar(jsonConfig, "config", "", "json config")

	flag.Parse()

	if *cfg.JSONConfig != "" {
		jsonConfigData, err := os.ReadFile(*cfg.JSONConfig)
		if err != nil {
			return err
		}

		var jsonValues map[string]any
		err = json.Unmarshal(jsonConfigData, &jsonValues)
		if err != nil {
			return err
		}

		for cfgKey, cfgValue := range jsonValues {
			val, ok := cfgValue.(string)
			if !ok {
				continue
			}

			switch cfgKey {
			case "address":
				if *cfg.Addr == defaultAgentAddr {
					*cfg.Addr = val
				}
			case "crypto_key":
				if *cfg.CryptoKey == defaultAgentCryptoKey {
					*cfg.CryptoKey = val
				}
			case "report_interval":
				if *cfg.ReportInterval == defaultReportInterval {
					interval, convErr := convStrIntervalToInt(val)
					if convErr != nil {
						return convErr
					} else if interval != 0 {
						*cfg.ReportInterval = interval
					}
				}
			case "poll_interval":
				if *cfg.PollInterval == defaultPoolInterval {
					interval, convErr := convStrIntervalToInt(val)
					if convErr != nil {
						return convErr
					} else if interval != 0 {
						*cfg.PollInterval = interval
					}
				}
			}
		}
	}

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
