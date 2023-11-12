// Package config for configure services
// Package will be parsing configs
package config

import (
	"flag"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

// Config - standard config interface
type Config interface {
	EnvParse() error
	ConfigureFlags()
}

type AddressConfig struct {
	Addr *string `env:"ADDRESS"`
}

type ReportConfig struct {
	AddressConfig

	ReportInterval *int `env:"REPORT_INTERVAL"`
	PollInterval   *int `env:"POLL_INTERVAL"`
}

type StoreConfig struct {
	StoreInterval   *uint   `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
}

type DBConfig struct {
	Dsn *string `env:"DATABASE_DSN"`
}

type HashKeyConfig struct {
	Key *string `env:"KEY"`
}

type RateConfig struct {
	Limit *int `env:"RATE_LIMIT"`
}

func NewAddressConfig() *AddressConfig {
	return &AddressConfig{}
}

func NewStoreConfig() *StoreConfig {
	return &StoreConfig{}
}

func NewReportConfig() *ReportConfig {
	return &ReportConfig{}
}

func NewDBConfig() *DBConfig {
	return &DBConfig{}
}

func NewHashKeyConfig() *HashKeyConfig {
	return &HashKeyConfig{}
}

func NewRateConfig() *RateConfig {
	return &RateConfig{}
}

// ParseCfgs - parse any configs
func ParseCfgs(cfgs []Config) error {
	for _, cfg := range cfgs {
		if err := cfg.EnvParse(); err != nil {
			return err
		}

		cfg.ConfigureFlags()
	}

	flag.Parse()

	return nil
}

func (cfg *AddressConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *AddressConfig) GetNormalizedAddress() string {
	if nil == cfg.Addr {
		return ""
	}

	if utf8.RuneCountInString(*cfg.Addr) < 4 || (*cfg.Addr)[:4] != "http" {
		return "http://" + *cfg.Addr
	}

	return *cfg.Addr
}

func (cfg *AddressConfig) ConfigureFlags() {
	addr := new(string)
	if cfg.Addr == nil {
		cfg.Addr = addr
	}

	flag.StringVar(addr, "a", ":8080", "http server addr")
}

func (cfg *StoreConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *StoreConfig) ConfigureFlags() {
	storeInterval := new(uint)
	fileStoragePath := new(string)
	restore := new(bool)

	if cfg.StoreInterval == nil {
		cfg.StoreInterval = storeInterval
	}

	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = fileStoragePath
	}

	if cfg.Restore == nil {
		cfg.Restore = restore
	}

	flag.UintVar(storeInterval, "i", 1, "store interval")
	flag.StringVar(fileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(restore, "r", true, "restore")
}

func (cfg *ReportConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *ReportConfig) ConfigureFlags() {
	addr := new(string)
	reportInterval := new(int)
	pollInterval := new(int)

	if cfg.Addr == nil {
		cfg.Addr = addr
	}

	if cfg.ReportInterval == nil {
		cfg.ReportInterval = reportInterval
	}

	if cfg.PollInterval == nil {
		cfg.PollInterval = pollInterval
	}

	flag.StringVar(addr, "a", "http://localhost:8080", "http server addr")
	flag.IntVar(reportInterval, "r", 10, "report interval")
	flag.IntVar(pollInterval, "p", 2, "poll interval")
}

func (cfg *DBConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *DBConfig) ConfigureFlags() {
	dsn := new(string)

	if cfg.Dsn == nil {
		cfg.Dsn = dsn
	}

	flag.StringVar(dsn, "d", "", "database dsn")
}

func (cfg *HashKeyConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *HashKeyConfig) ConfigureFlags() {
	key := new(string)

	if cfg.Key == nil {
		cfg.Key = key
	}

	flag.StringVar(key, "k", "", "hash key")
}

func (cfg *RateConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *RateConfig) ConfigureFlags() {
	limit := new(int)

	if cfg.Limit == nil {
		cfg.Limit = limit
	}

	flag.IntVar(limit, "l", 0, "rate limit")
}
