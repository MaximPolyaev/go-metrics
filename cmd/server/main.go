package main

import (
	"github.com/MaximPolyaev/go-metrics/internal/db"
	"github.com/MaximPolyaev/go-metrics/internal/storage/filestorage"
	"github.com/MaximPolyaev/go-metrics/internal/storage/memstorage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/router"
	"github.com/MaximPolyaev/go-metrics/internal/services/metricservice"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewAddressConfig()
	storeCfg := config.NewStoreConfig()
	dbConfig := config.NewDBConfig()

	err := config.ParseCfgs([]config.Config{
		cfg,
		storeCfg,
		dbConfig,
	})

	if err != nil {
		return err
	}

	lg := logger.New(os.Stdout)
	metricService, err := metricservice.New(
		memstorage.New(),
		filestorage.New(*storeCfg.FileStoragePath),
		storeCfg,
		lg,
	)
	if err != nil {
		return err
	}

	dbConn, err := db.InitDB(*dbConfig.Dsn)
	if err != nil {
		return err
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			lg.Error(err)
		}
	}()

	h := handler.New(metricService)

	shutdownHandler(metricService)

	return http.ListenAndServe(
		cfg.GetNormalizedAddress(),
		router.CreateRouter(h, lg, dbConn),
	)
}

func shutdownHandler(s *metricservice.MetricService) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		s.Sync()

		os.Exit(0)
	}()
}
