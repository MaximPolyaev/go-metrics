package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"unicode/utf8"
)

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
	if cfg.Addr == nil {
		cfg.Addr = new(string)
		*cfg.Addr = ":8080"
	}

	flag.StringVar(cfg.Addr, "a", *cfg.Addr, "http server addr")
}

func (cfg *StoreConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *StoreConfig) ConfigureFlags() {
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = new(uint)
		*cfg.StoreInterval = 1
	}

	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = new(string)
		*cfg.FileStoragePath = "/tmp/metrics-db.json"
	}

	if cfg.Restore == nil {
		cfg.Restore = new(bool)
		*cfg.Restore = true
	}

	flag.UintVar(cfg.StoreInterval, "i", *cfg.StoreInterval, "store interval")
	flag.StringVar(cfg.FileStoragePath, "f", *cfg.FileStoragePath, "file storage path")
	flag.BoolVar(cfg.Restore, "r", *cfg.Restore, "restore")
}

func (cfg *ReportConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *ReportConfig) ConfigureFlags() {
	if cfg.Addr == nil {
		cfg.Addr = new(string)
		*cfg.Addr = "http://localhost:8080"
	}

	if cfg.ReportInterval == nil {
		cfg.ReportInterval = new(int)
		*cfg.ReportInterval = 10
	}

	if cfg.PollInterval == nil {
		cfg.PollInterval = new(int)
		*cfg.PollInterval = 2
	}

	flag.StringVar(cfg.Addr, "a", *cfg.Addr, "http server addr")
	flag.IntVar(cfg.ReportInterval, "r", *cfg.ReportInterval, "report interval")
	flag.IntVar(cfg.PollInterval, "p", *cfg.PollInterval, "poll interval")
}

func (cfg *DBConfig) EnvParse() error {
	return env.Parse(cfg)
}

func (cfg *DBConfig) ConfigureFlags() {
	if cfg.Dsn == nil {
		cfg.Dsn = new(string)
	}

	flag.StringVar(cfg.Dsn, "d", *cfg.Dsn, "database dsn")
}
