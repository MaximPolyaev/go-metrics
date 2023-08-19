package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	wantURLReqs := []string{
		"/update/counter/PollCount/0/",
		"/update/gauge/TotalAlloc/0.000000/",
		"/update/gauge/HeapInuse/0.000000/",
		"/update/gauge/HeapSys/0.000000/",
		"/update/gauge/StackSys/0.000000/",
		"/update/gauge/NumForcedGC/0.000000/",
		"/update/gauge/NumGC/0.000000/",
		"/update/gauge/MSpanInuse/0.000000/",
		"/update/gauge/MSpanSys/0.000000/",
		"/update/gauge/NextGC/0.000000/",
		"/update/gauge/HeapIdle/0.000000/",
		"/update/gauge/Mallocs/0.000000/",
		"/update/gauge/Frees/0.000000/",
		"/update/gauge/PauseTotalNs/0.000000/",
		"/update/gauge/HeapObjects/0.000000/",
		"/update/gauge/HeapReleased/0.000000/",
		"/update/gauge/MCacheInuse/0.000000/",
		"/update/gauge/Lookups/0.000000/",
		"/update/gauge/OtherSys/0.000000/",
		"/update/gauge/StackInuse/0.000000/",
		"/update/gauge/Sys/0.000000/",
		"/update/gauge/Alloc/0.000000/",
		"/update/gauge/GCCPUFraction/0.000000/",
		"/update/gauge/HeapAlloc/0.000000/",
		"/update/gauge/BuckHashSys/0.000000/",
		"/update/gauge/LastGC/0.000000/",
		"/update/gauge/MCacheSys/0.000000/",
		"/update/gauge/GCSys/0.000000/",
		"/update/gauge/RandomValue/0.000000/",
	}

	var urlReqs []string

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlReqs = append(urlReqs, r.URL.Path)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewHTTPClient(srv.URL)

	stats := metric.Stats{}

	err := client.UpdateMetrics(&stats)
	assert.NoError(t, err)

	assert.ElementsMatch(t, wantURLReqs, urlReqs)
}
