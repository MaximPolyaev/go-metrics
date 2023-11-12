package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/handler"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/MaximPolyaev/go-metrics/internal/router"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	srv, err := makeTestHttpServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	req := prepareTestMetricsReq(srv.URL)

	do, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, do.StatusCode)
}

func BenchmarkUpdateMetrics(b *testing.B) {
	srv, err := makeTestHttpServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()

	req := prepareTestMetricsReq(srv.URL)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = http.DefaultClient.Do(req)
	}
}

func makeTestHttpServer() (*httptest.Server, error) {
	lg := logger.New(os.Stdout)
	storeCfg := makeTestStoreConfig()

	metricService, err := initMetricService(nil, storeCfg, nil)
	if err != nil {
		return nil, err
	}

	h := handler.New(metricService)
	r := router.CreateRouter(h, lg, nil, nil)

	return httptest.NewServer(r), nil
}

func makeTestStoreConfig() *config.StoreConfig {
	fileStoragePath := "/tmp/test-metrics.json"
	restore := true
	storeInterval := uint(1)

	storeCfg := config.NewStoreConfig()
	storeCfg.FileStoragePath = &fileStoragePath
	storeCfg.Restore = &restore
	storeCfg.StoreInterval = &storeInterval

	return storeCfg
}

func prepareTestMetricsReq(baseURL string) *http.Request {
	const maxMetricsOnType = 1000
	mSlice := make([]metric.Metric, 0, maxMetricsOnType)

	for i := 0; i < maxMetricsOnType; i++ {
		iStr := strconv.Itoa(i)

		gauge := metric.Metric{
			ID:    "gauge" + iStr,
			MType: metric.GaugeType,
			Value: new(float64),
		}
		*gauge.Value = float64(i)

		counter := metric.Metric{
			ID:    "gauge" + iStr,
			MType: metric.CounterType,
			Delta: new(int64),
		}
		*counter.Delta = int64(i)

		mSlice = append(mSlice, gauge, counter)
	}

	body, _ := json.Marshal(mSlice)
	buf := bytes.NewBuffer(body)

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/updates", buf)
	request.Header.Add("Content-Type", "application/json")

	return request
}
