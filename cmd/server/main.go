package main

import (
	"log"
	"net/http"
	"os"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/memstorage"
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
	if err := cfg.Parse(); err != nil {
		return err
	}

	store := memstorage.New()
	metricService := metricservice.New(store)

	h := handler.New(metricService)
	lg := logger.New(os.Stdout)

	muxRouter := router.CreateRouter(h, lg)

	return http.ListenAndServe(*cfg.Addr, muxRouter)
}
