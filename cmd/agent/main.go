package main

import (
	"log"
	"net/url"
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

	addr, err := url.Parse(*cfg.Addr)
	if err != nil {
		return err
	}

	httpClient := httpclient.NewHTTPClient(addr.String())

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
