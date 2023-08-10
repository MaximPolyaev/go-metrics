package server

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/config"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
)

func Run() error {
	cfg := config.NewServer()
	if err := cfg.Parse(); err != nil {
		return err
	}

	s := memstorage.NewMemStorage()

	muxRouter := router.CreateRouter(s)

	return http.ListenAndServe(*cfg.Addr, muxRouter)
}
