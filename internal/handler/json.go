package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

// BatchUpdateByJSONFunc - batch update metrics by json
func (h *Handler) BatchUpdateByJSONFunc() http.HandlerFunc {
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
		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				http.Error(w, closeErr.Error(), http.StatusInternalServerError)
			}
		}()

		var mSlice []metric.Metric
		if unmarshalErr := json.Unmarshal(buf, &mSlice); unmarshalErr != nil {
			http.Error(w, unmarshalErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		err = h.metricService.BatchUpdate(r.Context(), mSlice)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
}

// UpdateByJSONFunc - update metric by json
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
		}
		err = r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mm := new(metric.Metric)

		if unmarshalErr := json.Unmarshal(buf, mm); unmarshalErr != nil {
			http.Error(w, unmarshalErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		if validateErr := mm.ValidateWithValue(); validateErr != nil {
			http.Error(w, validateErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		mm = h.metricService.Update(r.Context(), mm)

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

// GetValueByJSONFunc - get metric value by json
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

		if unmarshalErr := json.Unmarshal(buf, mm); unmarshalErr != nil {
			http.Error(w, unmarshalErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		if validateErr := mm.Validate(); validateErr != nil {
			http.Error(w, validateErr.Error(), http.StatusUnprocessableEntity)
			return
		}

		mm, ok := h.metricService.Get(r.Context(), mm)
		if !ok {
			mm.ValueInit()
		}

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
