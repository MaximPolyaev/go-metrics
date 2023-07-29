package server

import (
	"net/http"
)

func incorrectMetricHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Incorrect metric type", http.StatusBadRequest)
}

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	gaugeMetric, err := makeGaugeMetricByUrlPath(r.URL.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	getMemStorage().writeGaugeMetric(*gaugeMetric)
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	counterMetric, err := makeCounterMetricByUrlPath(r.URL.Path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	getMemStorage().writeCounterMetric(*counterMetric)
}
