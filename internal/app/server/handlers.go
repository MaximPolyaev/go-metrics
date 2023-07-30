package server

import (
	"net/http"
	"strings"
)

func incorrectMetricHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Incorrect metric type", http.StatusBadRequest)
}

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	urlPath = strings.Trim(urlPath, "/")

	if len(urlPath) == 0 {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	gaugeMetric, err := makeGaugeMetricByUrlPath(r.URL.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getMemStorage().writeGaugeMetric(*gaugeMetric)
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	urlPath = strings.Trim(urlPath, "/")

	if len(urlPath) == 0 {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	counterMetric, err := makeCounterMetricByUrlPath(urlPath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	getMemStorage().writeCounterMetric(*counterMetric)
}
