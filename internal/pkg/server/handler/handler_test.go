package handler_test

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/encoding"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/memstorage"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateFunc(t *testing.T) {
	tests := []struct {
		name         string
		URL          string
		expectedCode int
	}{
		{
			name:         "gauge case #1",
			URL:          "/update/gauge/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #2",
			URL:          "/update/gauge/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #3",
			URL:          "/update/gauge/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "gauge case #4",
			URL:          "/update/gauge/test/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #5",
			URL:          "/update/gauge/test/2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "counter case #1",
			URL:          "/update/counter/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #2",
			URL:          "/update/counter/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #3",
			URL:          "/update/counter/test/test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "counter case #4",
			URL:          "/update/counter/test/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #5",
			URL:          "/update/counter/test/2",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.URL, nil)
			w := httptest.NewRecorder()

			muxRouter := router.CreateRouter(memstorage.NewMemStorage())
			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestMainFunc(t *testing.T) {
	storage := memstorage.NewMemStorage()

	binaryF, err := encoding.Float64ToByte(1.1)
	assert.NoError(t, err)

	binaryI := encoding.IntToByte(1)

	storage.Set(string(metric.GaugeType), "test", binaryF)
	storage.Set(string(metric.CounterType), "test", binaryI)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	muxRouter := router.CreateRouter(storage)
	muxRouter.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	str := w.Body.String()

	assert.Contains(t, str, "<ul><li>test: 1.1</li><li>test: 1</li></ul>")
}

func TestGetValue(t *testing.T) {
	tests := []struct {
		name            string
		URL             string
		expectedCode    int
		expectedBodyStr string
	}{
		{
			name:         "gauge case #1",
			URL:          "/value/gauge/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #2",
			URL:          "/value/gauge/notExist",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #3",
			URL:          "/value/gauge/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:            "gauge case #4",
			URL:             "/value/gauge/test",
			expectedCode:    http.StatusOK,
			expectedBodyStr: "1.1",
		},
		{
			name:            "gauge case #5",
			URL:             "/value/gauge/test/",
			expectedCode:    http.StatusOK,
			expectedBodyStr: "1.1",
		},
		{
			name:         "counter case #1",
			URL:          "/value/counter/",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #2",
			URL:          "/value/counter/notExist",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #3",
			URL:          "/value/counter/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:            "counter case #4",
			URL:             "/value/counter/test",
			expectedCode:    http.StatusOK,
			expectedBodyStr: "1",
		},
		{
			name:            "counter case #5",
			URL:             "/value/counter/test/",
			expectedCode:    http.StatusOK,
			expectedBodyStr: "1",
		},
	}

	binaryF, err := encoding.Float64ToByte(1.1)
	assert.NoError(t, err)

	binaryInt := encoding.IntToByte(1)

	storage := memstorage.NewMemStorage()
	storage.Set(string(metric.GaugeType), "test", binaryF)
	storage.Set(string(metric.CounterType), "test", binaryInt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.URL, nil)
			w := httptest.NewRecorder()

			muxRouter := router.CreateRouter(storage)
			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code == http.StatusOK {
				assert.Equal(t, tt.expectedBodyStr, w.Body.String())
			}
		})
	}
}
