package main

import (
	"log"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/handler"
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
	muxRouter := router.CreateRouter(h)

	return http.ListenAndServe(*cfg.Addr, muxRouter)
}
