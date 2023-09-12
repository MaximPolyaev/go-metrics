package config

import (
	"flag"
	"unicode/utf8"

	"github.com/caarlos0/env/v9"
)

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

func (cfg *StoreConfig) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.StoreInterval == nil {
		cfg.StoreInterval = new(uint)
		flag.UintVar(cfg.StoreInterval, "i", 1, "store interval")
	}

	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = new(string)
		flag.StringVar(cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	}

	if cfg.Restore == nil {
		cfg.Restore = new(bool)
		flag.BoolVar(cfg.Restore, "r", true, "restore")
	}

	flag.Parse()

	return nil
}

func (cfg *ReportConfig) Parse() error {
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

func (cfg *DBConfig) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	if cfg.Dsn == nil {
		cfg.Dsn = new(string)
		flag.StringVar(cfg.Dsn, "d", "", "database dsn")

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
