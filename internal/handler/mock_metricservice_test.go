package handler_test

import (
	"github.com/MaximPolyaev/go-metrics/internal/metric"
)

type mockMetricService struct{}

func (s *mockMetricService) Update(_ *metric.Metric) *metric.Metric {
	return nil
}

func (s *mockMetricService) Get(mm *metric.Metric) (*metric.Metric, bool) {
	mSlice := s.GetAll()

	for _, mFromSlice := range mSlice {
		if mFromSlice.ID == mm.ID && mFromSlice.MType == mm.MType {
			return &mFromSlice, true
		}
	}

	return mm, false
}

func (s *mockMetricService) GetAll() []metric.Metric {
	var delta int64
	var value float64

	delta = 10
	value = 1.1

	return []metric.Metric{
		{
			ID:    "test",
			MType: metric.CounterType,
			Delta: &delta,
		},
		{
			ID:    "test",
			MType: metric.GaugeType,
			Value: &value,
		},
	}
}
