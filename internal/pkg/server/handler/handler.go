package handler

import (
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/html"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	metricService metricService
}

type metricService interface {
	Update(mType metric.Type, name string, value string) error
	GetValues(mType metric.Type) (map[string]string, error)
	GetValue(mType metric.Type, name string) (value string, ok bool, err error)
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

		_, err = io.WriteString(w, htmlDocument.AsString())

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
}

func (h *Handler) UpdateFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := metric.Type(chi.URLParam(r, "type"))
		name := chi.URLParam(r, "name")
		valueStr := chi.URLParam(r, "value")

		if err := h.metricService.Update(mType, name, valueStr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
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
