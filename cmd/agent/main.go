package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/config"
	"github.com/MaximPolyaev/go-metrics/internal/crypto"
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

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	if err := printAppInfo(); err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func printAppInfo() error {
	_, err := fmt.Println(
		"Build version:", buildVersion,
		"Build date:", buildDate,
		"Build commit:", buildCommit,
	)

	return err
}

func run() error {
	cfg := config.NewAgentConfig()
	if err := cfg.Parse(); err != nil {
		return err
	}

	cryptoEncoder, err := makeCryptoEncoder(*cfg.CryptoKey)
	if err != nil {
		return err
	}

	lg := logger.New(os.Stdout)

	mStats := defaultstats.New()
	gopStats := gopsutilstats.New(lg)

	httpClient := httpclient.NewHTTPClient(
		cfg.GetNormalizedAddress(),
		*cfg.HashKey,
		cryptoEncoder,
	)

	chRead := make(chan Stats)
	chReport := make(chan Stats)

	poolInterval := time.NewTicker(time.Duration(*cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(*cfg.ReportInterval) * time.Second)

	for w := 0; w < maxWorkerCount; w++ {
		go readStats(chRead)
	}

	pushRate := computePushWorkerCount(*cfg.RateLimit)

	for w := 0; w < pushRate; w++ {
		go updateMetrics(httpClient, chReport, lg)
	}

	ctx, cancel := context.WithCancel(context.Background())

	shutdownHandler(cancel)

	for {
		select {
		case <-poolInterval.C:
			pushToChannel(chRead, mStats, gopStats)
		case <-reportInterval.C:
			pushToChannel(chReport, mStats, gopStats)
		case <-ctx.Done():
			close(chRead)
			close(chReport)
		}
	}
}

func makeCryptoEncoder(cryptoKey string) (*crypto.Encoder, error) {
	if cryptoKey == "" {
		return nil, nil
	}

	publicKey, err := crypto.LoadPublicKey(cryptoKey)

	if err != nil {
		return nil, err
	}

	return crypto.NewCryptoEncoder(publicKey), nil
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

func pushToChannel(ch chan<- Stats, sts ...Stats) {
	for _, s := range sts {
		ch <- s
	}
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

func shutdownHandler(cancelFunc context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		cancelFunc()

		os.Exit(0)
	}()
}
