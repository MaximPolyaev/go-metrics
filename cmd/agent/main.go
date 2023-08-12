package main

import (
	"log"
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

	httpClient := httpclient.NewHTTPClient(*cfg.Addr)

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-poolInterval.C:
			metric.ReadStats(&mStats)
		case <-reportInterval.C:
			if err := httpClient.UpdateMetrics(&mStats); err != nil {
				log.Println(err)
			}
		}
	}
}
