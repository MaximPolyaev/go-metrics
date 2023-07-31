package handler

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIncorrectMetric(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	IncorrectMetric(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGaugeFunc(t *testing.T) {
	tests := []struct {
		name         string
		stripUrl     string
		expectedCode int
	}{
		{
			name:         "case #1",
			stripUrl:     "/update/gauge/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "case #2",
			stripUrl:     "/update/gauge/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #3",
			stripUrl:     "/update/gauge/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #4",
			stripUrl:     "/update/gauge/test/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #5",
			stripUrl:     "/update/gauge/test/2",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.stripUrl, nil)
			w := httptest.NewRecorder()

			handler := http.StripPrefix("/update/gauge/", GaugeFunc(memstorage.NewMemStorage()))
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestCounterFunc(t *testing.T) {
	tests := []struct {
		name         string
		stripUrl     string
		expectedCode int
	}{
		{
			name:         "case #1",
			stripUrl:     "/update/counter/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "case #2",
			stripUrl:     "/update/counter/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #3",
			stripUrl:     "/update/counter/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #4",
			stripUrl:     "/update/counter/test/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "case #5",
			stripUrl:     "/update/counter/test/2",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.stripUrl, nil)
			w := httptest.NewRecorder()

			handler := http.StripPrefix("/update/counter/", CounterFunc(memstorage.NewMemStorage()))
			handler.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
