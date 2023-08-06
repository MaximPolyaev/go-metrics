package server

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/env"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/flags"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
	"net/http"
)

func Run() error {
	e, err := env.ParseEnv()
	if err != nil {
		return err
	}

	f := flags.ParseFlags(e)

	s := memstorage.NewMemStorage()

	muxRouter := router.CreateRouter(s)

	return http.ListenAndServe(f.GetAddr(), muxRouter)
}
