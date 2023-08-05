package server

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
	"net/http"
)

func Run(addr string) error {
	s := memstorage.NewMemStorage()

	muxRouter := router.CreateRouter(s)

	return http.ListenAndServe(addr, muxRouter)
}
