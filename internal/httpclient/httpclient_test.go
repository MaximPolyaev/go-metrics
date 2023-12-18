package httpclient

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaximPolyaev/go-metrics/internal/stats/defaultstats"

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

		reader, err := gzip.NewReader(r.Body)
		assert.NoError(t, err)

		buf, err := io.ReadAll(reader)
		assert.NoError(t, err)

		var records []struct {
			ID string `json:"id"`
		}

		err = json.Unmarshal(buf, &records)
		assert.NoError(t, err)

		for _, r := range records {
			idsFromReqs = append(idsFromReqs, r.ID)
		}

		hash := r.Header.Get("HashSHA256")

		assert.True(t, hash != "")
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	hashKey := "hash_key"
	client := NewHTTPClient(srv.URL, hashKey, nil)

	stats := defaultstats.New()

	err := client.UpdateMetrics(stats.AsMetrics())
	assert.NoError(t, err)

	assert.ElementsMatch(t, wantURLReqs, idsFromReqs)
}
