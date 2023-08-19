package router

import (
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	updateAction        = "/update/"
	valueAction         = "/value/"
	updateMetricPattern = updateAction + "{type}/{name}/{value}"
	getMetricPattern    = valueAction + "{type}/{name}"
)

type handler interface {
	UpdateFunc() http.HandlerFunc
	GetValueFunc() http.HandlerFunc
	MainFunc() http.HandlerFunc
}

func CreateRouter(h handler, log *logger.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.WithLogging(log))

	router.Post(updateMetricPattern, h.UpdateFunc())
	router.Post(updateMetricPattern+"/", h.UpdateFunc())
	router.Get(getMetricPattern, h.GetValueFunc())
	router.Get(getMetricPattern+"/", h.GetValueFunc())

	router.Get("/", h.MainFunc())

	return router
}
