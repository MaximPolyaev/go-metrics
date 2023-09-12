package main

import (
	"log"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewReportConfig()

	if err := config.ParseCfgs([]config.Config{cfg}); err != nil {
		return err
	}

	var mStats metric.Stats

	httpClient := httpclient.NewHTTPClient(cfg.GetNormalizedAddress())

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	for {
		select {
		case <-poolInterval.C:
			metric.ReadStats(&mStats)
		case <-reportInterval.C:
			if err := httpClient.UpdateMetrics(mStats.AsMetrics()); err != nil {
				log.Println(err)
			}
		}
	}
}
