package main

import (
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/server/handler"
	"github.com/MaximPolyaev/go-metrics/internal/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/server/router"
	"github.com/MaximPolyaev/go-metrics/internal/server/services/metricservice"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.NewServer()
	if err := cfg.Parse(); err != nil {
		return err
	}

	store := memstorage.New()
	metricService := metricservice.New(store)
	h := handler.New(metricService)
	muxRouter := router.CreateRouter(h)

	return http.ListenAndServe(*cfg.Addr, muxRouter)
}
