package server

import (
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/config"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/handler"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/services/metricservice"
)

func Run() error {
	cfg := config.NewServer()
	if err := cfg.Parse(); err != nil {
		return err
	}

	store := memstorage.New()
	metricService := metricservice.New(store)

	handlers := handler.New(&metricService)

	muxRouter := router.CreateRouter(&handlers)

	return http.ListenAndServe(*cfg.Addr, muxRouter)
}
