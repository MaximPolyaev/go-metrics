package router

import (
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const (
	updatePattern       = "/update"
	valuePattern        = "/value"
	updateMetricPattern = updatePattern + "/{type}/{name}/{value}"
	getMetricPattern    = valuePattern + "/{type}/{name}"
)

type handler interface {
	UpdateFunc() http.HandlerFunc
	GetValueFunc() http.HandlerFunc
	MainFunc() http.HandlerFunc
	UpdateByJSONFunc() http.HandlerFunc
}

func CreateRouter(h handler, log *logger.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.WithLogging(log))
	router.Use(chimiddleware.RedirectSlashes)

	router.Post(updatePattern, h.UpdateByJSONFunc())
	router.Post(updateMetricPattern, h.UpdateFunc())
	router.Get(getMetricPattern, h.GetValueFunc())

	router.Get("/", h.MainFunc())

	return router
}
