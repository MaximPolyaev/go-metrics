package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/MaximPolyaev/go-metrics/internal/html"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	metricService metricService
}

type metricService interface {
	GetValues(mType metric.Type) (map[string]string, error)
	GetValue(mType metric.Type, name string) (value string, ok bool, err error)
	Update(mm *metric.Metrics) *metric.Metrics
	Get(mm *metric.Metrics) *metric.Metrics
}

func New(mService metricService) *Handler {
	return &Handler{
		metricService: mService,
	}
}

func (h *Handler) MainFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list string

		values, err := h.metricService.GetValues(metric.GaugeType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for k, v := range values {
			list += html.Li(k + ": " + v)
		}

		values, err = h.metricService.GetValues(metric.CounterType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for k, v := range values {
			list += html.Li(k + ": " + v)
		}

		htmlDocument := html.NewDocument()
		htmlDocument.SetBody(html.Ul(list))

		w.Header().Set("Content-Type", "text/html")

		_, err = io.WriteString(w, htmlDocument.AsString())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
}

func (h *Handler) UpdateFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mm := metric.Metrics{
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

		h.metricService.Update(&mm)
	}
}

func (h *Handler) GetValueFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := metric.Type(chi.URLParam(r, "type"))
		name := chi.URLParam(r, "name")

		metricValue, ok, err := h.metricService.GetValue(mType, name)

		if err != nil {
			if ok {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			return
		}

		_, err = io.WriteString(w, metricValue)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
}
