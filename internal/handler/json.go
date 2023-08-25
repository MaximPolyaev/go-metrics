package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

func (h *Handler) UpdateByJSONFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(
				w,
				fmt.Sprintf("unexpected Content-Type: %s", ct),
				http.StatusBadRequest,
			)
			return
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mm := new(metric.Metric)

		if err := json.Unmarshal(buf, mm); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err := mm.ValidateWithValue(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		mm = h.metricService.Update(mm)

		resp, err := json.Marshal(mm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) GetValueByJSONFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(
				w,
				fmt.Sprintf("unexpected Content-Type: %s", ct),
				http.StatusBadRequest,
			)
			return
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mm := new(metric.Metric)

		if err := json.Unmarshal(buf, mm); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err := mm.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		mm = h.metricService.Get(mm)

		resp, err := json.Marshal(mm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
