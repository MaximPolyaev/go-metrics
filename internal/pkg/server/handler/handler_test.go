package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/handler"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/router"
	"github.com/stretchr/testify/assert"
)

type mockMetricService struct{}

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
	muxRouter := router.CreateRouter(h)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.URL, nil)
			w := httptest.NewRecorder()

			muxRouter.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestMainFunc(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h := handler.New(&mockMetricService{})
	muxRouter := router.CreateRouter(h)
	muxRouter.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	str := w.Body.String()

	assert.Contains(t, str, "<ul><li>test: 1.1</li><li>test: 10</li></ul>")
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
	muxRouter := router.CreateRouter(h)

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

func (m *mockMetricService) Update(_ metric.Type, _ string, _ string) error {
	return nil
}

func (m *mockMetricService) GetValues(mType metric.Type) (map[string]string, error) {
	switch mType {
	case metric.CounterType:
		return map[string]string{
			"test": "10",
		}, nil
	case metric.GaugeType:
		return map[string]string{
			"test": "1.1",
		}, nil
	}
	return nil, nil
}

func (m *mockMetricService) GetValue(mType metric.Type, name string) (value string, ok bool, err error) {
	if name == "notExist" {
		return "", false, errors.New("")
	}

	switch mType {
	case metric.CounterType:
		if name == "test" {
			return "10", true, nil
		}
	case metric.GaugeType:
		if name == "test" {
			return "1.1", true, nil
		}
	}

	return "", false, nil
}
