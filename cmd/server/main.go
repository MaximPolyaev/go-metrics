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
	printAppInfo()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func printAppInfo() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
}

func run() error {
	cfg := config.NewAddressConfig()
	storeCfg := config.NewStoreConfig()
	dbConfig := config.NewDBConfig()
	hashCfg := config.NewHashKeyConfig()
	cryptoCfg := config.NewCryptoConfig()

	configs := []config.Config{
		cfg,
		storeCfg,
		dbConfig,
		hashCfg,
		cryptoCfg,
	}
	err := config.ParseCfgs(configs)

	if err != nil {
		return err
	}

	lg := logger.New(os.Stdout)

	marshal, err := json.Marshal(configs)
	if err != nil {
		return err
	}

	lg.Info("configs ", string(marshal))

	cryptoDecoder, err := makeCryptoDecoder(cryptoCfg)
	if err != nil {
		return err
	}

	var dbConn *sql.DB

	if *dbConfig.Dsn != "" {
		dbConn, err = db.InitDB(*dbConfig.Dsn)
		if err != nil {
			return err
		}

		defer func() {
			if closeErr := dbConn.Close(); closeErr != nil {
				lg.Error(closeErr)
			}
		}()
	}

	metricService, err := initMetricService(
		dbConn,
		storeCfg,
		lg,
	)
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
		router.CreateRouter(h, lg, dbConn, hashCfg.Key, cryptoDecoder),
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

func makeCryptoDecoder(cryptoCfg *config.CryptoConfig) (*crypto.Decoder, error) {
	if cryptoCfg.CryptoKey == nil {
		return nil, nil
	}

	privateKey, err := crypto.LoadPrivateKey(*cryptoCfg.CryptoKey)
	if err != nil {
		return nil, err
	}

	return crypto.NewCryptoDecoder(privateKey), nil
}

func initMetricService(
	dbConn *sql.DB,
	storeCfg *config.StoreConfig,
	lg *logger.Logger,
) (*metricservice.MetricService, error) {
	if dbConn == nil {
		mService, err := metricservice.New(
			memstorage.New(),
			filestorage.New(*storeCfg.FileStoragePath),
			storeCfg,
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
