package mem_storage

import (
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"
	"github.com/MaximPolyaev/go-metrics/internal/pkg/server/storage"
)

type gaugeMetrics map[metric.Name]metric.Gauge
type counterMetrics map[metric.Name]metric.Counter

type MemStorage struct {
	gaugeMetrics   gaugeMetrics
	counterMetrics counterMetrics
}

func NewMemStorage() storage.Storage {
	return MemStorage{
		gaugeMetrics:   make(gaugeMetrics),
		counterMetrics: make(counterMetrics),
	}
}

func (s MemStorage) UpdateGaugeMetric(m metric.Gauge) {
	s.gaugeMetrics[m.Name] = m
}

func (s MemStorage) UpdateCounterMetric(m metric.Counter) {
	existMetric, ok := s.counterMetrics[m.Name]

	if ok {
		existMetric.Value += m.Value
		s.counterMetrics[m.Name] = existMetric
		return
	}

	s.counterMetrics[m.Name] = m
}
