package handler_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/MaximPolyaev/go-metrics/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateByJSONFunc(t *testing.T) {
	type Body struct {
		mm    metric.Metric
		delta int64
		value float64
	}

	tests := []struct {
		name        string
		method      string
		URL         string
		contentType string
		body        *Body
		wantStatus  int
	}{
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			URL:        "/update",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			URL:        "/update/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "empty content type",
			method:     http.MethodPost,
			URL:        "/update",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			URL:         "/update",
			contentType: "application/json",
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			URL:         "/update",
			contentType: "application/json",
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name:        "empty json fields",
			method:      http.MethodPost,
			URL:         "/update",
			contentType: "application/json",
			body: &Body{
				mm: metric.Metric{
					ID:    "",
					MType: metric.Type(""),
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "success counter",
			method:      http.MethodPost,
			URL:         "/update",
			contentType: "application/json",
			body: &Body{
				mm: metric.Metric{
					ID:    "test",
					MType: metric.CounterType,
				},
				delta: 10,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success gauge",
			method:      http.MethodPost,
			URL:         "/update",
			contentType: "application/json",
			body: &Body{
				mm: metric.Metric{
					ID:    "test",
					MType: metric.GaugeType,
				},
				value: 10,
			},
			wantStatus: http.StatusOK,
		},
	}

	h := handler.New(&mockMetricService{})
	lg := logger.New(os.Stdout)
	muxRouter := router.CreateRouter(h, lg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader io.Reader

			if tt.body != nil {
				tt.body.mm.Delta = new(int64)
				*tt.body.mm.Delta = tt.body.delta

				tt.body.mm.Value = new(float64)
				*tt.body.mm.Value = tt.body.value

				body, err := json.Marshal(tt.body.mm)
				fmt.Println(string(body))
				reader = strings.NewReader(string(body))
				assert.NoError(t, err)
			}

			r := httptest.NewRequest(tt.method, tt.URL, reader)
			w := httptest.NewRecorder()

			if tt.contentType != "" {
				r.Header.Set("Content-Type", tt.contentType)
			}

			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandler_GetValueByJSONFunc(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		URL         string
		contentType string
		mm          *metric.Metric
		wantStatus  int
	}{
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			URL:        "/value",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			URL:        "/value/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "empty content type",
			method:     http.MethodPost,
			URL:        "/value",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			URL:         "/value",
			contentType: "application/json",
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name:        "empty body",
			method:      http.MethodPost,
			URL:         "/value",
			contentType: "application/json",
			wantStatus:  http.StatusUnprocessableEntity,
		},
		{
			name:        "empty json fields",
			method:      http.MethodPost,
			URL:         "/value",
			contentType: "application/json",
			mm: &metric.Metric{
				ID:    "",
				MType: metric.Type(""),
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "success counter",
			method:      http.MethodPost,
			URL:         "/value",
			contentType: "application/json",
			mm: &metric.Metric{
				ID:    "test",
				MType: metric.CounterType,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "success gauge",
			method:      http.MethodPost,
			URL:         "/value",
			contentType: "application/json",
			mm: &metric.Metric{
				ID:    "test",
				MType: metric.GaugeType,
			},
			wantStatus: http.StatusOK,
		},
	}

	h := handler.New(&mockMetricService{})
	lg := logger.New(os.Stdout)
	muxRouter := router.CreateRouter(h, lg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader io.Reader

			if tt.mm != nil {
				body, err := json.Marshal(tt.mm)
				fmt.Println(string(body))
				reader = strings.NewReader(string(body))
				assert.NoError(t, err)
			}

			r := httptest.NewRequest(tt.method, tt.URL, reader)
			w := httptest.NewRecorder()

			if tt.contentType != "" {
				r.Header.Set("Content-Type", tt.contentType)
			}

			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
