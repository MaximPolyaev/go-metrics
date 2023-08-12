package router

import (
	"net/http"

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

func CreateRouter(h handler) *chi.Mux {
	router := chi.NewRouter()

	router.Post(updateMetricPattern, h.UpdateFunc())
	router.Post(updateMetricPattern+"/", h.UpdateFunc())
	router.Get(getMetricPattern, h.GetValueFunc())
	router.Get(getMetricPattern+"/", h.GetValueFunc())

	router.Get("/", h.MainFunc())

	return router
}
