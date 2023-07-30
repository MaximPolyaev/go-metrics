package handler

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/storage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/url_parser"
	"net/http"
	"strings"
)

func IncorrectMetric(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Incorrect metric type", http.StatusBadRequest)
}

func GaugeFunc(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.Trim(r.URL.Path, "/")

		if len(urlPath) == 0 {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		gaugeMetric, err := url_parser.MakeGaugeMetricByURLPath(urlPath)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.UpdateGaugeMetric(*gaugeMetric)
	}
}

func CounterFunc(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := strings.Trim(r.URL.Path, "/")

		if len(urlPath) == 0 {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		counterMetric, err := url_parser.MakeCounterMetricByURLPath(urlPath)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.UpdateCounterMetric(*counterMetric)
	}
}
