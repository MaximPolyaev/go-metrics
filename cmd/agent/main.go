package main

import (
	"log"
	"os"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/httpclient"
	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/MaximPolyaev/go-metrics/internal/metric"
	"github.com/MaximPolyaev/go-metrics/internal/stats/defaultstats"
	"github.com/MaximPolyaev/go-metrics/internal/stats/gopsutilstats"
)

const maxWorkerCount = 2
const minWorkerCount = 1

type Stats interface {
	ReadStats()
	AsMetrics() []metric.Metric
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewReportConfig()
	hashCfg := config.NewHashKeyConfig()
	rateCfg := config.NewRateConfig()

	if err := config.ParseCfgs([]config.Config{cfg, hashCfg, rateCfg}); err != nil {
		return err
	}

	lg := logger.New(os.Stdout)

	mStats := defaultstats.New()
	gopStats := gopsutilstats.New(lg)

	httpClient := httpclient.NewHTTPClient(cfg.GetNormalizedAddress(), hashCfg.Key)

	chRead := make(chan Stats)
	chPush := make(chan Stats)

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	for w := 0; w < maxWorkerCount; w++ {
		go readStats(chRead)
	}

	pushRate := computePushWorkerCount(*rateCfg.Limit)

	for w := 0; w < pushRate; w++ {
		go updateMetrics(httpClient, chPush, lg)
	}

	for {
		select {
		case <-poolInterval.C:
			chRead <- mStats
			chRead <- gopStats
		case <-reportInterval.C:
			chPush <- mStats
			chPush <- gopStats
		}
	}
}

func computePushWorkerCount(rateLimit int) int {
	if rateLimit > maxWorkerCount {
		return maxWorkerCount
	}

	if rateLimit < minWorkerCount {
		return minWorkerCount
	}

	return rateLimit
}

func readStats(chS <-chan Stats) {
	for s := range chS {
		s.ReadStats()
	}
}

func updateMetrics(httpClient *httpclient.HTTPClient, chS <-chan Stats, lg *logger.Logger) {
	for s := range chS {
		if err := httpClient.UpdateMetrics(s.AsMetrics()); err != nil {
			lg.Errorln(err)
		}
	}
}
