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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewBaseConfig()
	if err := cfg.Parse(); err != nil {
		return err
	}

	var mStats metric.Stats
	var wg sync.WaitGroup

	httpClient := httpclient.NewHTTPClient(*cfg.Addr)

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			<-poolInterval.C

			metric.ReadStats(&mStats)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			<-reportInterval.C

			if err := httpClient.UpdateMetrics(&mStats); err != nil {
				log.Println(err)
			}
		}
	}()

	wg.Wait()

	return nil
}
