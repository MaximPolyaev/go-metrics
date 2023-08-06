package agent

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/flags"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/metric"
	"sync"
	"time"
)

func Run() {
	f := flags.ParseFlags()

	var mStats metric.Stats
	var wg sync.WaitGroup
	wg.Add(1)

	httpClient := httpclient.NewHTTPClient(f.GetAddr())

	poolInterval := time.Duration(f.GetPollInterval()) * time.Second
	reportInterval := time.Duration(f.GetReportInterval()) * time.Second

	go func() {
		for {
			<-time.After(poolInterval)

			metric.ReadStats(&mStats)
		}
	}()

	go func() {
		for {
			<-time.After(reportInterval)

			if err := httpClient.UpdateMetrics(&mStats); err != nil {
				panic(err)
			}
		}
	}()

	wg.Wait()
}
