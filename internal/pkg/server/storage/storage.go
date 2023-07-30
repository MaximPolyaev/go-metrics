package storage

import "github.com/MaximPolyaev/go-metrics/internal/pkg/server/metric"

type Storage interface {
	UpdateGaugeMetric(m metric.Gauge)
	UpdateCounterMetric(m metric.Counter)
}
