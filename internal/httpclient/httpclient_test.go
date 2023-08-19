package httpclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetrics(t *testing.T) {
	wantURLReqs := []string{
		"PollCount",
		"TotalAlloc",
		"HeapInuse",
		"HeapSys",
		"StackSys",
		"NumForcedGC",
		"NumGC",
		"MSpanInuse",
		"MSpanSys",
		"NextGC",
		"HeapIdle",
		"Mallocs",
		"Frees",
		"PauseTotalNs",
		"HeapObjects",
		"HeapReleased",
		"MCacheInuse",
		"Lookups",
		"OtherSys",
		"StackInuse",
		"Sys",
		"Alloc",
		"GCCPUFraction",
		"HeapAlloc",
		"BuckHashSys",
		"LastGC",
		"MCacheSys",
		"GCSys",
		"RandomValue",
	}

	var idsFromReqs []string

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := r.Body.Close()
			assert.NoError(t, err)
		}()

		buf, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var record = struct {
			ID string `json:"id"`
		}{}

		err = json.Unmarshal(buf, &record)
		assert.NoError(t, err)

		idsFromReqs = append(idsFromReqs, record.ID)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	client := NewHTTPClient(srv.URL)

	stats := metric.Stats{}

	err := client.UpdateMetrics(stats.AsMetrics())
	assert.NoError(t, err)

	assert.ElementsMatch(t, wantURLReqs, idsFromReqs)
}
