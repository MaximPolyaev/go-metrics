package config

import (
	"encoding/json"
	"flag"
	"net"
	"os"

	"github.com/caarlos0/env/v9"
)

const (
	defaultServerAddr      = ":8080"
	defaultServerCryptoKey = ""
	defaultStoreInterval   = 1
	defaultStoreFile       = "/tmp/metrics-db.json"
	defaultDBDsn           = ""
	defaultRestore         = true
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
	JSONConfig      *string `env:"CONFIG"`
	TrustedSubnet   *Subnet `env:"TRUSTED_SUBNET"`
}

type Subnet net.IPNet

func (s *Subnet) UnmarshalText(text []byte) error {
	return s.UnmarshalString(string(text))
}

func (s *Subnet) UnmarshalString(str string) error {
	_, subnet, parseErr := net.ParseCIDR(str)
	if parseErr != nil {
		return parseErr
	}

	*s = Subnet(*subnet)
	return nil
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
	jsonConfig := new(string)
	trustedSubnet := new(string)

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
	if cfg.JSONConfig == nil {
		cfg.JSONConfig = jsonConfig
	}

	flag.StringVar(addr, "a", defaultServerAddr, "http server addr")
	flag.UintVar(storeInterval, "i", defaultStoreInterval, "store interval")
	flag.StringVar(fileStoragePath, "f", defaultStoreFile, "file storage path")
	flag.BoolVar(restore, "r", defaultRestore, "restore")
	flag.StringVar(DBDsn, "d", defaultDBDsn, "database dsn")
	flag.StringVar(hashKey, "k", "", "hash key")
	flag.StringVar(cryptoKey, "crypto-key", defaultServerCryptoKey, "crypto key")
	flag.StringVar(jsonConfig, "c", "", "json config")
	flag.StringVar(jsonConfig, "config", "", "json config")
	flag.StringVar(trustedSubnet, "t", "", "trusted subnet")

	flag.Parse()

	if cfg.TrustedSubnet == nil && *trustedSubnet != "" {
		cfg.TrustedSubnet = new(Subnet)
		err := cfg.TrustedSubnet.UnmarshalString(*trustedSubnet)
		if err != nil {
			return err
		}
	}

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
			switch cfgKey {
			case "address":
				if val, ok := cfgValue.(string); ok && *cfg.Addr == defaultServerAddr {
					*cfg.Addr = val
				}
			case "restore":
				if val, ok := cfgValue.(bool); ok && *cfg.Restore {
					*cfg.Restore = val
				}
			case "crypto_key":
				if val, ok := cfgValue.(string); ok && *cfg.CryptoKey == defaultServerCryptoKey {
					*cfg.CryptoKey = val
				}
			case "store_file":
				if val, ok := cfgValue.(string); ok && *cfg.FileStoragePath == defaultStoreFile {
					*cfg.FileStoragePath = val
				}
			case "database_dsn":
				if val, ok := cfgValue.(string); ok && *cfg.DBDsn == defaultDBDsn {
					*cfg.DBDsn = val
				}
			case "store_interval":
				if val, ok := cfgValue.(string); ok && *cfg.StoreInterval == defaultStoreInterval {
					interval, convErr := convStrIntervalToInt(val)
					if convErr != nil {
						return convErr
					} else if uint(interval) != 0 {
						*cfg.StoreInterval = uint(interval)
					}
				}
			case "trusted_subnet":
				if val, ok := cfgValue.(string); ok && cfg.TrustedSubnet == nil {
					cfg.TrustedSubnet = new(Subnet)
					unmarshalErr := cfg.TrustedSubnet.UnmarshalString(val)
					if unmarshalErr != nil {
						return unmarshalErr
					}
				}
			}
		}
	}

	return nil
}

func (cfg *ServerConfig) TrustedSubnetAsIpNet() *net.IPNet {
	if cfg.TrustedSubnet == nil {
		return nil
	}
	return (*net.IPNet)(cfg.TrustedSubnet)
}
