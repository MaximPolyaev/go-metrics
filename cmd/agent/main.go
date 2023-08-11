package main

import (
	"context"
	"github.com/MaximPolyaev/go-metrics/internal/agent/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/agent/metric"
	"github.com/MaximPolyaev/go-metrics/internal/agent/services/stats"
	"github.com/MaximPolyaev/go-metrics/internal/config"
	"log"
	"sync"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	cfg := config.NewBaseConfig()
	if err := cfg.Parse(); err != nil {
		log.Fatalln(err)
	}

	var mStats metric.Stats
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	httpClient := httpclient.NewHTTPClient(*cfg.Addr)
	statsService := stats.New(ctx, httpClient)

	go statsService.Pool(*cfg.PollInterval, &mStats)

	go func() {
		for {
			err = statsService.Report(*cfg.ReportInterval, &mStats)

			if err != nil {
				wg.Done()
				break
			}
		}
	}()

	wg.Wait()
	cancel()

	return
}
