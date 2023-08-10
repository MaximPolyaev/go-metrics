package agent

import (
	"log"
	"sync"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/agent/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/config"
)

func Run() {
	cfg := config.NewAgent()
	if err := cfg.Parse(); err != nil {
		log.Fatalln(err)
	}

	var mStats metric.Stats
	var wg sync.WaitGroup
	wg.Add(1)

	httpClient := httpclient.NewHTTPClient(*cfg.Addr)

	poolInterval := time.Duration(*cfg.PollInterval) * time.Second
	reportInterval := time.Duration(*cfg.ReportInterval) * time.Second

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
				log.Fatalln(err)
			}
		}
	}()

	wg.Wait()
}
