package main

import (
	"log"
	"sync"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/agent/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/agent/metric"
	"github.com/MaximPolyaev/go-metrics/internal/config"
)

func main() {
	run()
}

func run() {
	cfg := config.NewBaseConfig()
	if err := cfg.Parse(); err != nil {
		log.Fatalln(err)
	}

	var mStats metric.Stats
	var wg sync.WaitGroup
	wg.Add(1)

	httpClient := httpclient.NewHTTPClient(*cfg.Addr)

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	go func() {
		for {
			<-poolInterval.C

			metric.ReadStats(&mStats)
		}
	}()

	go func() {
		for {
			<-reportInterval.C

			if err := httpClient.UpdateMetrics(&mStats); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	wg.Wait()
}
