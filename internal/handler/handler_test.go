package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateFunc(t *testing.T) {
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
			URL:          "/update/gauge/test/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "gauge case #4",
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
			URL:          "/update/counter/test/test/test",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "counter case #4",
			URL:          "/update/counter/test/2",
			expectedCode: http.StatusOK,
		},
	}

	h := handler.New(&mockMetricService{})
	lg := logger.New(os.Stdout)
	muxRouter := router.CreateRouter(h, lg, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.URL, nil)
			w := httptest.NewRecorder()

			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestHandler_MainFunc(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := handler.New(&mockMetricService{})
	lg := logger.New(os.Stdout)
	muxRouter := router.CreateRouter(h, lg, nil, nil)
	muxRouter.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	str := w.Body.String()

	assert.Contains(t, str, "<ul><li>test: 10</li><li>test: 1.1</li></ul>")
}

func TestHandler_GetValue(t *testing.T) {
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
			expectedBodyStr: "10",
		},
		{
			name:            "counter case #5",
			URL:             "/value/counter/test/",
			expectedCode:    http.StatusOK,
			expectedBodyStr: "10",
		},
	}

	h := handler.New(&mockMetricService{})
	lg := logger.New(os.Stdout)
	muxRouter := router.CreateRouter(h, lg, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.URL, nil)
			w := httptest.NewRecorder()

			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code == http.StatusOK {
				assert.Equal(t, tt.expectedBodyStr, w.Body.String())
			}
		})
	}
}
