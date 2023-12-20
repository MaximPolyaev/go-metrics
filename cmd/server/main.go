package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/crypto"
	"github.com/MaximPolyaev/go-metrics/internal/db"
	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/router"
	"github.com/MaximPolyaev/go-metrics/internal/services/metricservice"
	"github.com/MaximPolyaev/go-metrics/internal/storage/dbstorage"
	"github.com/MaximPolyaev/go-metrics/internal/storage/filestorage"
	"github.com/MaximPolyaev/go-metrics/internal/storage/memstorage"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	if err := printAppInfo(); err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func printAppInfo() error {
	_, err := fmt.Println(
		"Build version:", buildVersion,
		"\nBuild date:", buildDate,
		"\nBuild commit:", buildCommit,
	)

	return err
}

func run() error {
	cfg := config.NewServerConfig()
	err := cfg.Parse()
	if err != nil {
		return err
	}

	lg := logger.New(os.Stdout)

	var jsonConfigs []byte
	jsonConfigs, err = json.Marshal(cfg)
	if err != nil {
		return err
	}
	lg.Info("configs ", string(jsonConfigs))

	cryptoDecoder, err := makeCryptoDecoder(*cfg.CryptoKey)
	if err != nil {
		return err
	}

	var dbConn *sql.DB

	if *cfg.DBDsn != "" {
		dbConn, err = db.InitDB(*cfg.DBDsn)
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := dbConn.Close(); closeErr != nil {
				lg.Error(closeErr)
			}
		}()
	}

	metricService, err := initMetricService(dbConn, cfg, lg)
	if err != nil {
		return err
	}

	if dbConn == nil {
		shutdownHandler(metricService)
	}

	h := handler.New(metricService)

	lg.Info("Start server on ", *cfg.Addr)

	return http.ListenAndServe(
		*cfg.Addr,
		router.CreateRouter(h, lg, dbConn, *cfg.HashKey, cryptoDecoder),
	)
}

func shutdownHandler(s *metricservice.MetricService) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		s.Sync(context.Background())

		os.Exit(0)
	}()
}

func makeCryptoDecoder(cryptoKey string) (*crypto.Decoder, error) {
	if cryptoKey == "" {
		return nil, nil
	}

	privateKey, err := crypto.LoadPrivateKey(cryptoKey)
	if err != nil {
		return nil, err
	}

	return crypto.NewCryptoDecoder(privateKey), nil
}

func initMetricService(
	dbConn *sql.DB,
	serverCfg *config.ServerConfig,
	lg *logger.Logger,
) (*metricservice.MetricService, error) {
	if dbConn == nil {
		mService, err := metricservice.New(
			memstorage.New(),
			filestorage.New(*serverCfg.FileStoragePath),
			serverCfg,
			lg,
		)
		if err != nil {
			return nil, err
		}
		return mService, nil
	}

	dbStorage := dbstorage.New(dbConn, lg)
	if err := dbStorage.Init(); err != nil {
		return nil, err
	}

	return metricservice.New(dbStorage, nil, nil, lg)
}
