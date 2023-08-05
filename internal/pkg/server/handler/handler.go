package handler

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/html"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/services"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MainFunc(s memstorage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gaugeService, err := services.FactoryMetricService(metric.GaugeType, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		counterService, err := services.FactoryMetricService(metric.CounterType, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		list := ""

		for _, mService := range []services.MetricService{gaugeService, counterService} {
			values, err := mService.GetValues()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for k, v := range values {
				list += html.Li(k + ": " + v)
			}
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

func UpdateFunc(s memstorage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := metric.Type(chi.URLParam(r, "type"))
		updateService, err := services.FactoryMetricService(mType, s)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := chi.URLParam(r, "name")
		valueStr := chi.URLParam(r, "value")

		if err := updateService.Update(name, valueStr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func GetValue(s memstorage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := metric.Type(chi.URLParam(r, "type"))
		metricService, err := services.FactoryMetricService(mType, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := chi.URLParam(r, "name")

		metricValue, ok, err := metricService.GetValue(name)

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
