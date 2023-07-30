package agent

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/metric"
	"sync"
	"time"
)

func Run(sPoolInterval int, sReportInterval int, baseUrl string) {
	var mStats metric.Stats
	var wg sync.WaitGroup
	wg.Add(1)

	httpClient := httpclient.NewHttpClient(baseUrl)

	poolInterval := time.Duration(sPoolInterval) * time.Second
	reportInterval := time.Duration(sReportInterval) * time.Second

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
