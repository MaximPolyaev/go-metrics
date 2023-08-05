package router

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/handler"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/go-chi/chi/v5"
)

const (
	updateAction        = "/update/"
	valueAction         = "/value/"
	updateMetricPattern = updateAction + "{type}/{name}/{value}"
	getMetricPattern    = valueAction + "{type}/{name}"
)

func CreateRouter(s memstorage.MemStorage) *chi.Mux {
	router := chi.NewRouter()

	router.Post(updateMetricPattern, handler.UpdateFunc(s))
	router.Post(updateMetricPattern+"/", handler.UpdateFunc(s))
	router.Get(getMetricPattern, handler.GetValue(s))
	router.Get(getMetricPattern+"/", handler.GetValue(s))

	router.Get("/", handler.MainFunc(s))

	return router
}
