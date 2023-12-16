// Package router configure middlewares and request patterns for run any handlers
package router

import (
	"database/sql"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

const (
	updatePattern       = "/update"
	updatesPattern      = "/updates"
	valuePattern        = "/value"
	updateMetricPattern = updatePattern + "/{type}/{name}/{value}"
	getMetricPattern    = valuePattern + "/{type}/{name}"
	pingPattern         = "/ping"
)

type handler interface {
	UpdateFunc() http.HandlerFunc
	GetValueFunc() http.HandlerFunc
	MainFunc() http.HandlerFunc
	UpdateByJSONFunc() http.HandlerFunc
	BatchUpdateByJSONFunc() http.HandlerFunc
	GetValueByJSONFunc() http.HandlerFunc
	PingFunc(db *sql.DB) http.HandlerFunc
}

type CryptoDecoder interface {
	Decode(data []byte) ([]byte, error)
}

func CreateRouter(
	h handler,
	log *logger.Logger,
	db *sql.DB,
	hashKey *string,
	cryptoDecoder CryptoDecoder,
) *chi.Mux {
	router := chi.NewRouter()

	if hashKey != nil && *hashKey != "" {
		router.Use(middleware.WithHashing(*hashKey))
	}

	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.WithLogging(log))
	router.Use(middleware.WithDecrypt(cryptoDecoder, updatesPattern+"/"))
	router.Use(chimiddleware.StripSlashes)

	router.Post(updatePattern, h.UpdateByJSONFunc())
	router.Post(updatesPattern, h.BatchUpdateByJSONFunc())
	router.Post(valuePattern, h.GetValueByJSONFunc())
	router.Post(updateMetricPattern, h.UpdateFunc())
	router.Get(getMetricPattern, h.GetValueFunc())

	router.Get(pingPattern, h.PingFunc(db))

	router.Get("/", h.MainFunc())

	return router
}
