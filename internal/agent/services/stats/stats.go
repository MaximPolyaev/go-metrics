package stats

import (
	"context"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/agent/metric"
)

type StatsService struct {
	ctx        context.Context
	httpClient httpClient
}

type httpClient interface {
	UpdateMetrics(stats *metric.Stats) error
}

func New(ctx context.Context, client httpClient) StatsService {
	return StatsService{
		ctx:        ctx,
		httpClient: client,
	}
}

func (s *StatsService) Pool(interval int, stats *metric.Stats) {
	tickerInterval := time.NewTicker(time.Duration(interval) * time.Second)
	defer tickerInterval.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tickerInterval.C:
			metric.ReadStats(stats)
		}
	}
}

func (s *StatsService) Report(interval int, stats *metric.Stats) error {
	tickerInterval := time.NewTicker(time.Duration(interval) * time.Second)
	defer tickerInterval.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-tickerInterval.C:
			if err := s.httpClient.UpdateMetrics(stats); err != nil {
				return err
			}
		}
	}
}
