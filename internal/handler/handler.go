package handler

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/html"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	metricService metricService
}

type metricService interface {
	Update(ctx context.Context, mm *metric.Metric) *metric.Metric
	Get(ctx context.Context, mm *metric.Metric) (*metric.Metric, bool)
	GetAll(ctx context.Context) []metric.Metric
	BatchUpdate(ctx context.Context, mSlice []metric.Metric) ([]metric.Metric, error)
}

func New(mService metricService) *Handler {
	return &Handler{
		metricService: mService,
	}
}

func (h *Handler) PingFunc(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) MainFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list string
		for _, m := range h.metricService.GetAll(r.Context()) {
			list += html.Li(m.ID + ": " + m.GetValueAsStr())
		}

		htmlDocument := html.NewDocument()
		htmlDocument.SetBody(html.Ul(list))

		w.Header().Set("Content-Type", "text/html")

		if _, err := io.WriteString(w, htmlDocument.AsString()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
}

func (h *Handler) UpdateFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mm := metric.Metric{
			ID:    chi.URLParam(r, "name"),
			MType: metric.Type(chi.URLParam(r, "type")),
		}

		valueStr := chi.URLParam(r, "value")

		switch mm.MType {
		case metric.GaugeType:
			value, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mm.Value = &value
		case metric.CounterType:
			value, err := strconv.ParseInt(valueStr, 0, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mm.Delta = &value
		}

		if err := mm.ValidateWithValue(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.metricService.Update(r.Context(), &mm)
	}
}

func (h *Handler) GetValueFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := metric.Type(chi.URLParam(r, "type"))
		name := chi.URLParam(r, "name")

		m := metric.Metric{ID: name, MType: mType}

		if err := m.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		metricValue, ok := h.metricService.Get(r.Context(), &m)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		if _, err := io.WriteString(w, metricValue.GetValueAsStr()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
}
