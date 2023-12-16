package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

// ServerConfig предназначен для настройки сервера сбора метрик
type ServerConfig struct {
	Addr            *string `env:"ADDRESS"`
	StoreInterval   *uint   `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
	DBDsn           *string `env:"DATABASE_DSN"`
	HashKey         *string `env:"KEY"`
	CryptoKey       *string `env:"CRYPTO_KEY"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (cfg *ServerConfig) Parse() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}

	addr := new(string)
	storeInterval := new(uint)
	fileStoragePath := new(string)
	restore := new(bool)
	DBDsn := new(string)
	hashKey := new(string)
	cryptoKey := new(string)

	if cfg.Addr == nil {
		cfg.Addr = addr
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = storeInterval
	}
	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = fileStoragePath
	}
	if cfg.Restore == nil {
		cfg.Restore = restore
	}
	if cfg.DBDsn == nil {
		cfg.DBDsn = DBDsn
	}
	if cfg.HashKey == nil {
		cfg.HashKey = hashKey
	}
	if cfg.CryptoKey == nil {
		cfg.CryptoKey = cryptoKey
	}

	flag.StringVar(addr, "a", ":8080", "http server addr")
	flag.UintVar(storeInterval, "i", 1, "store interval")
	flag.StringVar(fileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.BoolVar(restore, "r", true, "restore")
	flag.StringVar(DBDsn, "d", "", "database dsn")
	flag.StringVar(hashKey, "k", "", "hash key")
	flag.StringVar(cryptoKey, "crypto-key", "", "crypto key")

	flag.Parse()
	return nil
}
